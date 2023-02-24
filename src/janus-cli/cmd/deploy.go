package cmd

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"

	"github.com/spf13/cobra"
)

var (
	agentName string
	agentPort int

	hostName string

	localIp net.IP
)

// DeployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an aca-py agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if hostName != "" {
			valid := agent_deploy.ValidateSSHHostName(hostName)
			if !valid {
				log.Fatal("hostname flag must be on user@ip format")
			}
		}

		deployAgent()
	},
}

func init() {
	var err error

	localIp, err = agent_deploy.GetOutboundIP()
	if err != nil {
		log.Fatal(err)
	}

	deployCmd.Flags().StringVarP(&hostName, "host-name", "H", "", "The hostname of the target host you want to deploy the agent. The required format is user@ip.")
	deployCmd.Flags().StringVar(&agentName, "agent-name", "", "The aca-py agent name. This flag is optional but recommended so the agent will be better named")
	deployCmd.Flags().IntVar(&agentPort, "agent-port", 0, "The aca-py agent port. This flag is required and agent-port+1 will be used for aca-py admin page")

	deployCmd.MarkFlagRequired("agent-port")
	deployCmd.MarkFlagRequired("host-name")

	rootCmd.AddCommand(deployCmd)
}

func deployAgent() {
	agent := agent_deploy.AgentInfo{
		Name:      agentName,
		AdminPort: strconv.Itoa(agentPort + 1),
		AgentPort: strconv.Itoa(agentPort),
	}

	// Generate seed
	agent.Seed = agent_deploy.GenerateSeed()
	log.Printf("Seed generated: %s\n", agent.Seed)

	// Register DID
	did, err := agent_deploy.RegisterDID(agent.Seed, "http://dev.bcovrin.vonx.io")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DiD registered: %s\n", did)

	// Define endpoint
	if hostName != "" {
		agent.Endpoint = strings.Split(hostName, "@")[1]
	} else {
		agent.Endpoint = localIp.String()
	}

	// log results
	parsedAgent, err := json.Marshal(agent)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Deploying agent: %s\n", parsedAgent)

	// Instantiate Agent
	err = agent_deploy.InstantiateAgent(agent, hostName)
	if err != nil {
		log.Fatal(err)
	}
}
