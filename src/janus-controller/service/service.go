package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/Instituto-Atlantico/janus/pkg/agents"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/local"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/remote"
	"github.com/ldej/go-acapy-client"
)

type Service struct {
	ServerClient *acapy.Client
	Agents       map[string]*acapy.Client
}

func (s *Service) Init() {
	//add deploy only if no agent is running
	err := local.DeployAgent("192.168.0.10")
	if err != nil {
		log.Fatal(err)
	}

	s.ServerClient = acapy.NewClient("http://192.168.0.10:8002")
	s.Agents = make(map[string]*acapy.Client)
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
		// add deploy only if no agent is running
		err = remote.DeployAgent(provisionBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error Deploying agent: ", err)
			return
		}

		ip := strings.Split(provisionBody.DeviceHostName, "@")[1]
		s.Agents[ip] = acapy.NewClient(fmt.Sprintf("http://%s:8002", ip))

		go func() {
			helper.TryUntilNoError(s.Agents[ip].Status, 600) //check if agent is already up and running
			log.Println("Changing invitation")
			agents.ChangeInvitations(s.ServerClient, s.Agents[ip])
		}()

		b, _ := json.Marshal(provisionBody)
		fmt.Fprint(w, string(b))
	})

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, reflect.ValueOf(s.Agents).MapKeys())
	})

	log.Println("Server listening on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
