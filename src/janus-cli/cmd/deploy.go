package cmd

import (
	"fmt"
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
		panic(err)
	}

	deployCmd.Flags().StringVar(&name, "agent-name", hostname, "The aca-py agent name. This flag is optional, and will be filled with the host name if blank")
	deployCmd.Flags().IntVar(&agentPort, "agent-port", 0, "The aca-py agent port. This flag is required and agent-port+1 will be used for aca-py admin page if its available")

	deployCmd.MarkFlagRequired("agent-port")
	rootCmd.AddCommand(deployCmd)
}

func deployAgent(agentName string, agentPort int) {
	// - Discover IP
	ip, err := agent_deploy.GetOutboundIP()
	if err != nil {
		fmt.Println(err)
	}

	endpoint := fmt.Sprintf("http://%s", ip.String())
	fmt.Printf("Actual IP: %s\n", endpoint)

	// - Set Ports
	adminPort := agentPort + 1

	fmt.Printf("Agent port: %v Admin port: %v\n", agentPort, adminPort)

	// - Generate seed
	seed := agent_deploy.GenerateSeed()

	fmt.Printf("Seed generated: %s\n", seed)

	// - Register DID
	did, err := agent_deploy.RegisterDID(seed, "http://dev.bcovrin.vonx.io")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("DiD registered: %s\n", did)

	// - Instatiate Agent
	err = agent_deploy.InstatiateAgent(seed, agentName, strconv.Itoa(adminPort), strconv.Itoa(agentPort), endpoint)

	if err != nil {
		fmt.Println(err)
	}
}
