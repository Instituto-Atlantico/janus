package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Instituto-Atlantico/janus/pkg/agents"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/local"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/remote"
	"github.com/ldej/go-acapy-client"
)

type Device struct {
	Client       *acapy.Client
	ConnectionID string
}

type Service struct {
	ServerClient     *acapy.Client
	Agents           map[string]*Device
	CredDefinitionId string
}

func (s *Service) Init() {
	err := local.DeployAgent("192.168.0.10")
	if err != nil {
		log.Fatal(err)
	}

	s.ServerClient = acapy.NewClient("http://192.168.0.10:8002")

	helper.TryUntilNoError(s.ServerClient.Status, 600)

	// create cred definition
	s.CredDefinitionId, err = agents.CreateCredDef(s.ServerClient, "EZpfyRHcXuohyTvbgsrg7S:2:janus-sensors:1.0")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("CredDefinitionID", s.CredDefinitionId)

	s.Agents = make(map[string]*Device)
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

		for _, sensorType := range provisionBody.Permissions {
			permission := acapy.CredentialPreviewAttributeV2{
				MimeType: "text/plain",
				Name:     sensorType,
				Value:    "1",
			}

			permissionList = append(permissionList, permission)
		}

		fmt.Println(permissionList)

		device := Device{}

		ip := strings.Split(provisionBody.DeviceHostName, "@")[1]
		device.Client = acapy.NewClient(fmt.Sprintf("http://%s:8002", ip))

		go func() {
			helper.TryUntilNoError(device.Client.Status, 600) //check if agent is already up and running
			log.Println("Changing invitation")

			invitationID, _, _ := agents.ChangeInvitations(s.ServerClient, device.Client)
			device.ConnectionID = invitationID

			fmt.Println(device)

			s.Agents[ip] = &device

			time.Sleep(2 * time.Second)

			agents.IssueCredential(s.ServerClient, s.CredDefinitionId, device.ConnectionID, permissionList)
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
