package cmd

import (
	"github.com/spf13/cobra"
)

var (
	agentName string
	agentPort int
)

// DeployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an aca-py agent",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	deployCmd.PersistentFlags().StringVar(&agentName, "agent-name", "", "The aca-py agent name. This flag is optional but recommended so the agent will be better named")
	deployCmd.PersistentFlags().IntVar(&agentPort, "agent-port", 0, "The aca-py agent port. This flag is required and agent-port+1 will be used for aca-py admin page")

	deployCmd.MarkPersistentFlagRequired("agent-port")

	rootCmd.AddCommand(deployCmd)
}
