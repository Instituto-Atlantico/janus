/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"

	"github.com/spf13/cobra"
)

var (
	name string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy called:", name)
		deployAgent()
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().StringVar(&name, "name", "", "The aca-py agent name")

	rootCmd.AddCommand(deployCmd)
}

func deployAgent() {
	// init variables
	ledger := "http://dev.bcovrin.vonx.io"

	// - Discover IP
	ip, err := agent_deploy.GetOutboundIP()
	if err != nil {
		fmt.Println(err)
	}

	endpoint := fmt.Sprintf("http://%s", ip.String())
	fmt.Printf("Actual IP: %s\n", endpoint)

	// - Provision Port // needs to be done manually because it needs to be open in firewall
	agentPort := "8000"
	adminPort := "8001"

	fmt.Printf("Admin port: %s\nAgent port: %s\n", adminPort, agentPort)

	// - Generate seed
	seed := agent_deploy.GenerateSeed()

	fmt.Printf("Seed generated: %s\n", seed)

	// - Register DID
	did, err := agent_deploy.RegisterDID(seed, ledger)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("DiD registered: %s\n", did)

	// - Instatiate Agent
	err = agent_deploy.InstatiateAgent(seed, "vitor", adminPort, agentPort, ip.String())

	if err != nil {
		fmt.Println(err)
	}
}
