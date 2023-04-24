package cmd

import (
	"github.com/spf13/cobra"
)

var (
	agentServicePort string
	agentApiPort     string
)

// DeployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an issuer or Holder janus agents",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) { cmd.Println(cmd.Usage()) },
}

func init() {
	deployCmd.PersistentFlags().StringVar(&agentApiPort, "agent-api-port", "8002", "The aca-py agent api port, used for the controller and external agent requests.")
	deployCmd.PersistentFlags().StringVar(&agentServicePort, "agent-service-port", "8001", "The aca-py agent service port, used internally for communication with the ledger")

	rootCmd.AddCommand(deployCmd)
}
