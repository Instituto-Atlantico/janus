package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
	"github.com/spf13/cobra"
)

var (
	agentIp        string
	controllerPort string
)

var issuerCmd = &cobra.Command{
	Use:   "issuer",
	Short: "deploy an issuer aca-py agent locally and a janus-controller attached to it",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		deployAgentLocally()
	},
}

func init() {
	deployCmd.AddCommand(issuerCmd)

	issuerCmd.Flags().StringVar(&agentIp, "agent-ip", "", "The device`s ip of the network with the holder devices")
	issuerCmd.Flags().StringVar(&controllerPort, "controller-port", "8081", "The port that the janus-controller will run")
	issuerCmd.MarkFlagRequired("agent-ip")
}

func deployAgentLocally() {
	agent := agent_deploy.ComposeInfo{
		Name:           "issuer",
		AdminPort:      agentApiPort,
		AgentPort:      agentServicePort,
		Endpoint:       fmt.Sprintf("http://%s", agentIp),
		ControllerPort: controllerPort,
	}

	// generate seed and did
	seed, did := agent_deploy.ProvideDid()
	log.Printf("Seed generated: %s\n", seed)
	log.Printf("DiD registered: %s\n", did)

	agent.Seed = seed

	// log agent
	parsedAgent, err := json.Marshal(agent)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Deploying agent: %s\n", parsedAgent)

	// Instantiate Agent
	err = agent_deploy.InstantiateAgent(agent, "", "issuer")
	if err != nil {
		log.Fatal(err)
	}
}
