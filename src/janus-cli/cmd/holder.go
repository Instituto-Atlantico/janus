package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/pkg/yaml_parser"
	"github.com/spf13/cobra"
)

var (
	hostName string
	file     string
)

var holderCmd = &cobra.Command{
	Use:   "holder",
	Short: "Deploy an aca-py agent remotely",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if file != "" { //many hosts
			deployMultipleAgents()
		} else if hostName != "" { //single host
			valid := helper.ValidateSSHHostName(hostName)
			if !valid {
				log.Fatal("hostname flag must be on user@ip format")
			}

			deploySingleAgent()
		} else {
			log.Fatal("Either hostname or filename flags must passed")
		}
	},
}

func init() {
	deployCmd.AddCommand(holderCmd)

	holderCmd.Flags().StringVarP(&file, "config-file", "F", "", "The yaml file for deploy of multiple hosts at once. Check https://github.com/Instituto-Atlantico/janus/demo/agents.yaml for an example.")
	holderCmd.Flags().StringVarP(&hostName, "host-name", "H", "", "The hostname of the target host you want to deploy the agent. The required format is user@ip.")
}

func deployMultipleAgents() {
	fmt.Println("call multiple")

	body := yaml_parser.ParseFile(file)

	var wg sync.WaitGroup

	for _, a := range body.Agents {

		wg.Add(1)

		agent := agent_deploy.ComposeInfo{
			Name:      "holder",
			HostName:  a.Hostname,
			AdminPort: "8002",
			AgentPort: "8001",
			Endpoint:  fmt.Sprintf("http://%s", strings.Split(a.Hostname, "@")[1]),
		}

		go func() {
			log.Printf("deploy agent on %s\n", a.Hostname)
			defer wg.Done()

			err := deployAgent(agent)
			if err != nil {
				log.Println(err)
				runtime.Goexit() //interrupts the current routine
			}
		}()
	}
	wg.Wait()
}

func deploySingleAgent() {
	agent := agent_deploy.ComposeInfo{
		Name:      "holder",
		HostName:  hostName,
		AdminPort: agentApiPort,
		AgentPort: agentServicePort,
		Endpoint:  fmt.Sprintf("http://%s", strings.Split(hostName, "@")[1]),
	}

	err := deployAgent(agent)
	if err != nil {
		log.Fatal(err)
	}
}

func deployAgent(agent agent_deploy.ComposeInfo) error {
	// generate seed and did
	seed, did := agent_deploy.ProvideDid()
	log.Printf("Seed generated: %s\n", seed)
	log.Printf("DiD registered: %s\n", did)

	agent.Seed = seed

	// log agent
	parsedAgent, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	log.Printf("Deploying agent: %s\n", parsedAgent)

	// Instantiate Agent
	err = agent_deploy.InstantiateAgent(agent, "holder")
	if err != nil {
		return err
	}

	return nil
}
