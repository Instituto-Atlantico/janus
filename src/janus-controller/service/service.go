package service

import (
	"encoding/json"
	"fmt"
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

	log "github.com/Instituto-Atlantico/janus/pkg/logger"
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

func (s *Service) Init(serverAgentIp string) {
	var err error

	schemaId := "EZpfyRHcXuohyTvbgsrg7S:2:janus-sensors:1.0"

	s.ServerClient = acapy.NewClient(fmt.Sprintf("http://%s:8002", serverAgentIp))
	helper.TryUntilNoError(s.ServerClient.Status, 600)

	// create cred definition
	s.CredDefinitionId, err = agents.GetCredDef(s.ServerClient, schemaId)

	if err != nil {
		if err.Error() == "empty" {
			s.CredDefinitionId, err = agents.CreateCredDef(s.ServerClient, schemaId)
			if err != nil {
				log.FatalLogger(err)
			}
		} else {
			log.FatalLogger(err)
		}
	}

	log.InfoLogger("CredDefinitionID: %s", s.CredDefinitionId)

	s.Agents = make(map[string]*Device)
}

func (s *Service) RunCollector(timeoutInSeconds int) {
	ticker := time.NewTicker(time.Duration(timeoutInSeconds) * time.Second)

	go func() {
		for range ticker.C {
			ips := reflect.ValueOf(s.Agents).MapKeys()
			if len(ips) > 0 {
				log.InfoLogger("Getting sensor data for...")
				agentIP := ips[0]
				agentClient := s.Agents[agentIP.String()]

				sensorData := sensors.CollectSensorData(agentIP.String(), "5000")

				validatedData := make(map[string]any)

				for name, value := range sensorData {
					log.InfoLogger("Agent %s: Device %s sensor has value %s", agentIP, name, value)

					// request presentation proof for name
					presentationRequest, _ := agents.CreateRequestPresentationForSensor(s.ServerClient, s.CredDefinitionId, agentClient.ConnectionID, name)

					time.Sleep(2 * time.Second)

					credential, err := agents.GetCredential(agentClient.Client, "cred_def_id", s.CredDefinitionId)
					if err != nil {
						//log.Println(err)

						return
					}

					agents.SendPresentationByID(agentClient.Client, presentationRequest, credential)

					// wait for presentation to be ready
					_, err = helper.TryUntilNoError(func() ([]acapy.PresentationExchangeRecord, error) {
						return agents.IsPresentationDone(s.ServerClient, presentationRequest.ThreadID)
					}, 20)
					if err != nil {
						log.InfoLogger("Timeout presentation done")

						return
					}

					// if presentation is valid store value
					result, err := agents.VerifyPresentationByID(s.ServerClient, presentationRequest)
					if err != nil {
						//log.Println(err)

						return
					}

					if result.Verified == "true" {
						validatedData[name] = value
					}

				}

				log.InfoLogger("Agent %s: Validating device sensors permissions...", agentIP)
				log.InfoLogger("Agent %s: Allowed sensor data %s", agentIP, validatedData)

				// send sensor data to Dojot upon presentation proof
				log.InfoLogger("Agent %s: Publishing message to Dojot...", agentIP)
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
		var provisionBody ProvisionBody
		err := json.NewDecoder(r.Body).Decode(&provisionBody)

		if !ProvisionBodyIsValid(provisionBody) || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid body")

			return
		}

		ip := strings.Split(provisionBody.DeviceHostName, "@")[1]

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

		device := Device{}

		device.Client = acapy.NewClient(fmt.Sprintf("http://%s:8002", ip))
		device.BrokerUsername = provisionBody.BrokerUsername

		s.Broker = mqtt_pub.BrokerData{
			BrokerServerIp: provisionBody.BrokerServerIp,
			BrokerPassword: provisionBody.BrokerPassword,
		}

		go func() {
			helper.TryUntilNoError(device.Client.Status, 30) //check if agent is already up and running

			log.InfoLogger("Agent %s: Exchanging invitation...", ip)

			invitationID, _, err := agents.ChangeInvitations(s.ServerClient, device.Client)
			device.ConnectionID = invitationID
			if err != nil {
				log.InfoLogger("Agent %s: %s", ip, err)
			}

			log.InfoLogger("Agent %s: Invitation accepted with ID %s", ip, invitationID)

			time.Sleep(5 * time.Second)

			log.InfoLogger("Agent %s: Getting credential...", ip)
			cred, err := agents.GetCredential(device.Client, "cred_def_id", s.CredDefinitionId) //issue new credential only if no previous created
			if err != nil && err.Error() == "empty" {
				log.InfoLogger("Agent %s: Issuing credential", ip)

				agents.IssueCredential(s.ServerClient, s.CredDefinitionId, device.ConnectionID, permissionList)
				cred, err = helper.TryUntilNoError(func() (acapy.Credential, error) {
					return agents.GetCredential(device.Client, "cred_def_id", s.CredDefinitionId)
				}, 20)

				if err != nil {
					log.InfoLogger("Timeout on agents.GetCredential")
				}
			}

			log.InfoLogger("Agent %s: Device credential %s ", ip, cred)

			s.Agents[ip] = &device
		}()

		b, _ := json.Marshal(provisionBody)
		fmt.Fprint(w, string(b))
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reflect.ValueOf(s.Agents).MapKeys())
	})

	log.InfoLogger("Server listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		//log.Println(err)
	}
}
