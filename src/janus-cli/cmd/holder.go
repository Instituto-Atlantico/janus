package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/spf13/cobra"
)

var (
	hostName string
)

var holderCmd = &cobra.Command{
	Use:   "holder",
	Short: "Deploy an aca-py agent remotely",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		valid := helper.ValidateSSHHostName(hostName)
		if !valid {
			log.Fatal("hostname flag must be on user@ip format")
		}

		deployAgentRemotely()
	},
}

func init() {
	deployCmd.AddCommand(holderCmd)

	holderCmd.Flags().StringVarP(&hostName, "host-name", "H", "", "The hostname of the target host you want to deploy the agent. The required format is user@ip.")
	holderCmd.MarkFlagRequired("host-name")
}

func deployAgentRemotely() {
	agent := agent_deploy.ComposeInfo{
		Name:      "holder",
		AdminPort: agentApiPort,
		AgentPort: agentServicePort,
		Endpoint:  fmt.Sprintf("http://%s", strings.Split(hostName, "@")[1]),
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
	err = agent_deploy.InstantiateAgent(agent, hostName, "holder")
	if err != nil {
		log.Fatal(err)
	}
}
