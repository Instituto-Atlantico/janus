/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/spf13/cobra"
)

var (
	agentEndpoint string
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "deploy an aca-py agent locally",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		deployAgentLocally()
	},
}

func init() {
	deployCmd.AddCommand(localCmd)

	localIp, err := helper.GetOutboundIP()
	if err != nil {
		log.Fatal(err)
	}

	localCmd.Flags().StringVar(&agentEndpoint, "agent-ip", localIp.String(), "The ip that other agents will use to communicate with this. If blank it will be taken automatically.")
}

func deployAgentLocally() {
	agent := agent_deploy.AgentInfo{
		Name:      agentName,
		AdminPort: strconv.Itoa(agentPort + 1),
		AgentPort: strconv.Itoa(agentPort),
		Endpoint:  fmt.Sprintf("http://%s", agentEndpoint),
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
	err = agent_deploy.InstantiateAgent(agent, "", "server", true)
	if err != nil {
		log.Fatal(err)
	}
}
