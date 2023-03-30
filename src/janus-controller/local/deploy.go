package local

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
)

func DeployAgent(ip string) error {
	// Instantiate Agent
	agent := agent_deploy.AgentInfo{
		Name:      "janus-issuer",
		AdminPort: strconv.Itoa(8002),
		AgentPort: strconv.Itoa(8001),
		Endpoint:  fmt.Sprintf("http://%s", ip),
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
	err = agent_deploy.InstantiateAgent(agent, "", "server", false)
	if err != nil {
		return err
	}

	return nil
}
