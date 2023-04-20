package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Instituto-Atlantico/go-acapy-client"
	"github.com/Instituto-Atlantico/janus/pkg/agents"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/pkg/mqtt_pub"
	"github.com/Instituto-Atlantico/janus/pkg/sensors"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/local"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/remote"
)

type Device struct {
	Client         *acapy.Client
	ConnectionID   string
	BrokerUsername string
}

type Service struct {
	ServerClient     *acapy.Client
	Agents           map[string]*Device
	CredDefinitionId string
	Broker           mqtt_pub.BrokerData
}

var AllowedPermissions = []string{
	"temperature", "humidity",
}

func (s *Service) Init() {
	schemaId := "EZpfyRHcXuohyTvbgsrg7S:2:janus-sensors:1.0"

	err := local.DeployAgent("192.168.0.5")
	if err != nil {
		log.Fatal(err)
	}

	s.ServerClient = acapy.NewClient("http://192.168.0.5:8002")

	helper.TryUntilNoError(s.ServerClient.Status, 600)

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
				publicationTopic := fmt.Sprintf("%s/attrs", agentClient.BrokerUsername)
				mqtt_pub.PublishMessage(s.Broker, agentClient.BrokerUsername, publicationTopic, validatedData)
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
		var provisionBody remote.ProvisionBody
		err := json.NewDecoder(r.Body).Decode(&provisionBody)

		if !remote.ProvisionBodyIsValid(provisionBody) || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid body")

			return
		}

		//deploy agent
		err = remote.DeployAgent(provisionBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error Deploying agent: ", err)
			return
		}

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

		device := Device{}

		ip := strings.Split(provisionBody.DeviceHostName, "@")[1]
		device.Client = acapy.NewClient(fmt.Sprintf("http://%s:8002", ip))
		device.BrokerUsername = provisionBody.BrokerUsername

		s.Broker = mqtt_pub.BrokerData{
			BrokerServerIp: provisionBody.BrokerServerIp,
			BrokerPassword: provisionBody.BrokerPassword,
		}

		go func() {
			helper.TryUntilNoError(device.Client.Status, 600) //check if agent is already up and running
			log.Println("\nChanging invitation for agent ", ip)

			invitationID, _, _ := agents.ChangeInvitations(s.ServerClient, device.Client)
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
				}
			}

			log.Println("Device`s cred: ", cred)

			s.Agents[ip] = &device
		}()

		b, _ := json.Marshal(provisionBody)
		fmt.Fprint(w, string(b))
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reflect.ValueOf(s.Agents).MapKeys())
	})

	log.Println("Server listening on port ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}
