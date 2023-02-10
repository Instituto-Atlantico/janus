package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"

	"github.com/spf13/cobra"
)

var (
	name      string
	agentPort int
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an aca-py agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		deployAgent(name, int(agentPort))
	},
}

func init() {
	//gets host name to use if argument --name hasn't been passed
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	deployCmd.Flags().StringVar(&name, "agent-name", hostname, "The aca-py agent name. This flag is optional, and will be filled with the host name if blank")
	deployCmd.Flags().IntVar(&agentPort, "agent-port", 0, "The aca-py agent port. This flag is required and agent-port+1 will be used for aca-py admin page if its available")
	deployCmd.MarkFlagRequired("agent-port")

	rootCmd.AddCommand(deployCmd)
}

func getIP() (string, error) {
	ip, err := agent_deploy.GetOutboundIP()
	if err != nil {
		return "", err
	}

	formattedIP := fmt.Sprintf("http://%s", ip.String())
	return formattedIP, nil
}

func deployAgent(agentName string, agentPort int) {
	// - Set adminPort
	adminPort := agentPort + 1
	log.Printf("Agent port: %v Admin port: %v\n", agentPort, adminPort)

	// - Discover IP
	endpoint, err := getIP()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Actual IP: %s\n", endpoint)

	// - Generate seed
	seed := agent_deploy.GenerateSeed()
	log.Printf("Seed generated: %s\n", seed)

	// - Register DID
	did, err := agent_deploy.RegisterDID(seed, "http://dev.bcovrin.vonx.io")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DiD registered: %s\n", did)

	// - Instatiate Agent
	err = agent_deploy.InstatiateAgent(seed, agentName, strconv.Itoa(adminPort), strconv.Itoa(agentPort), endpoint)
	if err != nil {
		log.Fatal(err)
	}
}
