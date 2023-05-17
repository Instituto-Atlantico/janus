package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/Instituto-Atlantico/go-acapy-client"
	"github.com/Instituto-Atlantico/janus/pkg/agent_deploy"
	"github.com/Instituto-Atlantico/janus/pkg/controller_handlers"
	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/pkg/yaml_parser"
	"github.com/spf13/cobra"
)

var (
	hostName       string
	file           string
	provision      bool
	controllerAddr string
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
	holderCmd.Flags().BoolVarP(&provision, "provision", "p", false, "Allows auto-provision of the agent on the controller. The flag --controller-address specifies the controller ip  worand is setted as locahost:8081 as default")
	holderCmd.Flags().StringVarP(&controllerAddr, "controller", "c", "localhost:8081", "The controller address used by auto-provision process for the file deploy")
}

func deployMultipleAgents() {
	fmt.Println("call multiple")

	body := yaml_parser.ParseFile(file)

	var wg sync.WaitGroup

	for _, a := range body.Agents {

		wg.Add(1)

		composeInfo := agent_deploy.ComposeInfo{
			Name:      "holder",
			HostName:  a.Hostname,
			AdminPort: "8002",
			AgentPort: "8001",
			Endpoint:  fmt.Sprintf("http://%s", strings.Split(a.Hostname, "@")[1]),
		}

		go func(agent yaml_parser.Agent, provision bool) {
			defer wg.Done()

			log.Printf("deploying agent on %s device\n", agent.Hostname)
			err := deployAgent(composeInfo)
			if err != nil {
				log.Println(err)
				runtime.Goexit() //interrupts the current routine
			}

			if provision {
				ip := strings.Split(agent.Hostname, "@")[1]
				client := acapy.NewClient(fmt.Sprintf("http://%s:8002", ip))

				log.Printf("Waiting for agent deploy on %s device\n", agent.Hostname)
				_, err = helper.TryUntilNoError(func() (acapy.Status, error) {
					return client.Status()
				}, 600)
				if err != nil {
					log.Printf("error on auto-provisioning for agent %s:%s\n", ip, err)
				}

				service := controller_handlers.ControllerService{
					Address: controllerAddr,
				}

				body := controller_handlers.ProvisionBody{
					DeviceHostName: agent.Hostname,
					Permissions:    agent.Sensors,
					BrokerIp:       agent.Broker.IP,
					BrokerUsername: fmt.Sprintf("%s:%s", agent.Broker.Username, agent.Broker.ID),
					BrokerPassword: agent.Broker.Password,
				}
				log.Printf("Provisioning agent of %s device on %s controller\n", agent.Hostname, controllerAddr)
				err = service.RequestProvision(body)
				if err != nil {
					log.Println("Error on provisioning: ", err)
					runtime.Goexit() //interrupts the current routine
				}

				log.Println("Provisioning done for agent ", agent.Hostname)
			}
		}(a, provision)
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
