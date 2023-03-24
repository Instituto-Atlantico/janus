/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/spf13/cobra"
)

var (
	hostName string
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
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
	remoteCmd.Flags().StringVarP(&hostName, "host-name", "H", "", "The hostname of the target host you want to deploy the agent. The required format is user@ip.")
	remoteCmd.MarkFlagRequired("host-name")

	deployCmd.AddCommand(remoteCmd)
}

func deployAgentRemotely() {
	agent := agent_deploy.AgentInfo{
		Name:      agentName,
		AdminPort: strconv.Itoa(agentPort + 1),
		AgentPort: strconv.Itoa(agentPort),
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
	err = agent_deploy.InstantiateAgent(agent, hostName, "raspberry")
	if err != nil {
		log.Fatal(err)
	}
}
