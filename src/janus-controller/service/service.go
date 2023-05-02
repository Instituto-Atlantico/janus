package service

import (
	"encoding/json"
	"fmt"
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

	log "github.com/Instituto-Atlantico/janus/pkg/logger"
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
}

var AllowedPermissions = []string{
	"temperature", "humidity",
}

func (s *Service) Init(serverAgentIp string) {
	var err error

	schemaId := "EZpfyRHcXuohyTvbgsrg7S:2:janus-sensors:1.0"

	s.ServerClient = acapy.NewClient(fmt.Sprintf("http://%s:8002", serverAgentIp))
	_, err = helper.TryUntilNoError(s.ServerClient.Status, 20)
	if err != nil {
		log.FatalLogger("Timeout when trying to connect with issuer aca-py agent")
	}

	// create cred definition
	log.InfoLogger("Looking for a valid credential definition")

	s.CredDefinitionId, err = agents.GetCredDef(s.ServerClient, schemaId)

	if err != nil {
		if err.Error() == "empty" {
			log.InfoLogger("No previously created credential definition found. Issuing a new one")

			s.CredDefinitionId, err = agents.CreateCredDef(s.ServerClient, schemaId)
			if err != nil {
				log.FatalLogger(err)
			}
		} else {
			log.FatalLogger(err)
		}
	}

	log.InfoLogger("Credential definition created with ID: %s", s.CredDefinitionId)

	s.Agents = make(map[string]*Device)
}

func (s *Service) RunCollector(timeoutInSeconds int) {
	ticker := time.NewTicker(time.Duration(timeoutInSeconds) * time.Second)

	go func() {
		for range ticker.C {
			ips := reflect.ValueOf(s.Agents).MapKeys()
			if len(ips) > 0 {

				agentIP := ips[0] // pendind multidevice in parallel
				agentClient := s.Agents[agentIP.String()]
				log.InfoLogger("Agent %s: collecting sensor data", agentIP)

				sensorData, err := sensors.CollectSensorData(agentIP.String(), "5000")
				if err != nil {
					log.ErrorLogger("Agent %s: error collecting sensor data: %s", agentIP, err)
					continue
				}

				validatedData := make(map[string]any)

				for name, value := range sensorData {
					log.InfoLogger("Agent %s: Device %s sensor has value %s", agentIP, name, value)

					// request presentation proof for name
					presentationRequest, err := agents.CreateRequestPresentationForSensor(s.ServerClient, s.CredDefinitionId, agentClient.ConnectionID, name)
					if err != nil {
						log.ErrorLogger("Agent %s: error creating presentation request for sensor %s: %s", agentIP, name, err)
						continue
					}

					time.Sleep(2 * time.Second)

					credential, err := agents.GetCredential(agentClient.Client, "cred_def_id", s.CredDefinitionId)
					if err != nil {
						log.ErrorLogger("Agent %s: error getting device credential: %s", agentIP, err)
						continue
					}

					agents.SendPresentationByID(agentClient.Client, presentationRequest, credential)

					// wait for presentation to be ready
					_, err = helper.TryUntilNoError(func() ([]acapy.PresentationExchangeRecord, error) {
						return agents.IsPresentationDone(s.ServerClient, presentationRequest.ThreadID)
					}, 20)
					if err != nil {
						log.ErrorLogger("Agent %s: Timeout waiting for presentation done", agentIP)
						continue
					}

					// if presentation is valid store value
					log.InfoLogger("Agent %s: Validating device sensors permissions", agentIP)
					result, err := agents.VerifyPresentationByID(s.ServerClient, presentationRequest)
					if err != nil {
						log.ErrorLogger("Agent %s: error verifying presentation proof: %s", agentIP, err)
						continue
					}

					if result.Verified == "true" {
						validatedData[name] = value
					}

				}
				log.InfoLogger("Agent %s: Allowed sensor data %s", agentIP, validatedData)

				// send sensor data to Dojot upon presentation proof
				log.InfoLogger("Agent %s: Publishing message to Dojot", agentIP)
				err = mqtt_pub.PublishMessage(agentClient.BrokerCredentials, validatedData)
				if err != nil {
					log.ErrorLogger("Agent %s: error publishing message to mqtt Broker: %s", agentIP, err)
					continue
				}
				log.InfoLogger("Agent %s: message sent to Dojot", agentIP)
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
			Ip:       provisionBody.BrokerIp,
			Username: provisionBody.BrokerUsername,
			Password: provisionBody.BrokerPassword,
		}

		_, err = device.Client.Status()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.ErrorLogger("Device agent is not running properly")
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

		log.InfoLogger("Agent %s: Device permission list %s", ip, permissionList)

		b, _ := json.Marshal(provisionBody)
		fmt.Fprint(w, string(b))

		go func() {
			log.InfoLogger("Agent %s: Exchanging invitation...", ip)
			invitationID, _, err := agents.ChangeInvitations(s.ServerClient, device.Client)
			if err != nil {
				log.ErrorLogger("Agent %s: %s", ip, err)
				runtime.Goexit() //interrupts the current routine
			}
			device.ConnectionID = invitationID
			if err != nil {
				log.ErrorLogger("Agent %s: %s", ip, err)
				runtime.Goexit() //interrupts the current routine
			}

			log.InfoLogger("Agent %s: Invitation accepted with ID %s", ip, invitationID)

			time.Sleep(5 * time.Second)

			log.InfoLogger("Agent %s: Looking for a valid credential", ip)
			cred, err := agents.GetCredential(device.Client, "cred_def_id", s.CredDefinitionId) //issue new credential only if no previous created
			if err != nil && err.Error() == "empty" {
				log.InfoLogger("Agent %s: No previously created credential found. Issuing a new credential", ip)

				agents.IssueCredential(s.ServerClient, s.CredDefinitionId, device.ConnectionID, permissionList)
				cred, err = helper.TryUntilNoError(func() (acapy.Credential, error) {
					return agents.GetCredential(device.Client, "cred_def_id", s.CredDefinitionId)
				}, 20)
				if err != nil {
					log.ErrorLogger("Agent %s: Timeout getting a credential", ip)
					runtime.Goexit() //interrupts the current routine
				}
			}

			log.InfoLogger("Agent %s: Device credential %s ", ip, cred)

			s.Agents[ip] = &device
		}()
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reflect.ValueOf(s.Agents).MapKeys())
	})

	log.InfoLogger("Server listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.ErrorLogger("Server listening: %s", err)
	}
}
