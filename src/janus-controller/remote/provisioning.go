package remote

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
)

type ProvisionBody struct {
	DeviceHostName string   `json:"deviceHostName"` // user@ip
	Permissions    []string `json:"permissions"`    // ["temperature", "humidity"]
	BrokerServerIp string   `json:"brokerServerIp"` // brokerIP
	BrokerUsername string   `json:"brokerUsername"` // brokerUser:deviceID
	BrokerPassword string   `json:"brokerPassword"` // brokerPass
}

func ProvisionBodyIsValid(body ProvisionBody) bool {
	if body.DeviceHostName == "" || len(body.Permissions) == 0 || body.BrokerServerIp == "" || body.BrokerUsername == "" || body.BrokerPassword == "" {
		return false
	}

	return true
}

func DeployAgent(provisionBody ProvisionBody) error {
	// Instantiate Agent
	agent := agent_deploy.AgentInfo{
		Name:      "holder",
		AdminPort: strconv.Itoa(8002),
		AgentPort: strconv.Itoa(8001),
		Endpoint:  fmt.Sprintf("http://%s", strings.Split(provisionBody.DeviceHostName, "@")[1]),
	}

	// generate seed and did
	seed, did := agent_deploy.ProvideDid()
	log.Printf("Seed generated: %s\n", seed)
	log.Printf("DiD registered: %s\n", did)

	agent.Seed = seed

	// log agent
	parsedAgent, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	log.Printf("Deploying agent: %s\n", parsedAgent)

	// Instantiate Agent
	go agent_deploy.InstantiateAgent(agent, provisionBody.DeviceHostName, "raspberry", false)

	return nil
}
