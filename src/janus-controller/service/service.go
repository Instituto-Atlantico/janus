package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Instituto-Atlantico/go-acapy-client"
	"github.com/Instituto-Atlantico/janus/pkg/agents"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/pkg/mqtt_pub"
	"github.com/Instituto-Atlantico/janus/pkg/sensors"
)

type Device struct {
	Client            *acapy.Client
	ConnectionID      string
	BrokerCredentials mqtt_pub.BrokerCredentials
}

type Service struct {
	ServerClient     *acapy.Client
	Agents           map[string]*Device
	CredDefinitionId string
	BrokerIp         string
}

var AllowedPermissions = []string{
	"temperature", "humidity",
}

func (s *Service) Init(serverAgentIp, brokerIp string) {
	var err error

	schemaId := "EZpfyRHcXuohyTvbgsrg7S:2:janus-sensors:1.0"

	s.ServerClient = acapy.NewClient(fmt.Sprintf("http://%s:8002", serverAgentIp))
	_, err = helper.TryUntilNoError(s.ServerClient.Status, 20)
	if err != nil {
		log.Fatal("Timeout when trying to connect with issuer aca-py agent")
	}

	// create cred definition
	s.CredDefinitionId, err = agents.GetCredDef(s.ServerClient, schemaId)

	if err != nil {
		if err.Error() == "empty" {
			s.CredDefinitionId, err = agents.CreateCredDef(s.ServerClient, schemaId)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	log.Println("CredDefinitionID: ", s.CredDefinitionId)

	s.BrokerIp = brokerIp

	s.Agents = make(map[string]*Device)
}

func (s *Service) RunCollector(timeoutInSeconds int) {
	ticker := time.NewTicker(time.Duration(timeoutInSeconds) * time.Second)

	go func() {
		for range ticker.C {
			ips := reflect.ValueOf(s.Agents).MapKeys()
			if len(ips) > 0 {
				fmt.Println("Getting sensors data")
				agentIP := ips[0]
				agentClient := s.Agents[agentIP.String()]

				sensorData := sensors.CollectSensorData(agentIP.String(), "5000")

				validatedData := make(map[string]any)

				for name, value := range sensorData {
					fmt.Printf("Sensor [%s] has Value [%s]\n", name, value)

					// request presentation proof for name
					presentationRequest, _ := agents.CreateRequestPresentationForSensor(s.ServerClient, s.CredDefinitionId, agentClient.ConnectionID, name)

					time.Sleep(2 * time.Second)

					credential, err := agents.GetCredential(agentClient.Client, "cred_def_id", s.CredDefinitionId)
					if err != nil {
						log.Println(err)

						return
					}

					agents.SendPresentationByID(agentClient.Client, presentationRequest, credential)

					// wait for presentation to be ready
					_, err = helper.TryUntilNoError(func() ([]acapy.PresentationExchangeRecord, error) {
						return agents.IsPresentationDone(s.ServerClient, presentationRequest.ThreadID)
					}, 20)
					if err != nil {
						log.Println("Timeout presentation done")

						return
					}

					// if presentation is valid store value
					result, err := agents.VerifyPresentationByID(s.ServerClient, presentationRequest)
					if err != nil {
						log.Println(err)

						return
					}

					if result.Verified == "true" {
						validatedData[name] = value
					}
				}

				fmt.Println("Validate data:", validatedData)

				// send sensor data to Dojot upon presentation proof
				fmt.Println("Publishing message to Dojot...")
				mqtt_pub.PublishMessage(s.BrokerIp, agentClient.BrokerCredentials, validatedData)
			}
		}
	}()
}

func (s *Service) RunApi(port string) {
	http.HandleFunc("/provision", func(w http.ResponseWriter, r *http.Request) {
		//check method
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Invalid method")
			return
		}

		// parse body
		var provisionBody ProvisionBody
		err := json.NewDecoder(r.Body).Decode(&provisionBody)

		if !ProvisionBodyIsValid(provisionBody) || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid body")

			return
		}

		//create device object
		device := Device{}

		ip := strings.Split(provisionBody.DeviceHostName, "@")[1]
		device.Client = acapy.NewClient(fmt.Sprintf("http://%s:8002", ip))

		device.BrokerCredentials = mqtt_pub.BrokerCredentials{
			Username: provisionBody.BrokerUsername,
			Password: provisionBody.BrokerPassword,
		}

		_, err = device.Client.Status()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Device agent is not running properly")
			return
		}

		// Parse Permissions to credential previews
		permissionList := make([]acapy.CredentialPreviewAttributeV2, 0)

		for _, sensorType := range AllowedPermissions {

			allowed := helper.SliceContains(provisionBody.Permissions, sensorType)

			permission := acapy.CredentialPreviewAttributeV2{
				MimeType: "text/plain",
				Name:     sensorType,
				Value:    strconv.FormatBool(allowed),
			}

			permissionList = append(permissionList, permission)
		}

		fmt.Println(permissionList)

		b, _ := json.Marshal(provisionBody)
		fmt.Fprint(w, string(b))

		go func() {
			log.Println("\nChanging invitation for agent ", ip)
			invitationID, _, err := agents.ChangeInvitations(s.ServerClient, device.Client)
			if err != nil {
				log.Println(err)
				runtime.Goexit() //interrupts the current routine
			}
			device.ConnectionID = invitationID

			time.Sleep(5 * time.Second)

			cred, err := agents.GetCredential(device.Client, "cred_def_id", s.CredDefinitionId) //issue new credential only if no previous created
			if err != nil && err.Error() == "empty" {
				log.Println("\nIssuing credential for agent ", ip)
				agents.IssueCredential(s.ServerClient, s.CredDefinitionId, device.ConnectionID, permissionList)
				cred, err = helper.TryUntilNoError(func() (acapy.Credential, error) {
					return agents.GetCredential(device.Client, "cred_def_id", s.CredDefinitionId)
				}, 20)
				if err != nil {
					log.Println("Timeout on agents.GetCredential")
					runtime.Goexit() //interrupts the current routine
				}
			}

			log.Println("Device`s cred: ", cred)

			s.Agents[ip] = &device
		}()
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reflect.ValueOf(s.Agents).MapKeys())
	})

	log.Println("Server listening on port ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
