package agent_deploy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Instituto-Atlantico/janus/pkg/helper"
	"github.com/Instituto-Atlantico/janus/pkg/logger"
)

type ComposeInfo struct {
	Seed      string
	Name      string
	HostName  string
	Endpoint  string
	AdminPort string
	AgentPort string

	ControllerPort string
}

func parseComposeInfo(info ComposeInfo, cmd *exec.Cmd) {
	cmd.Env = os.Environ()

	//Aca-py
	cmd.Env = append(cmd.Env, fmt.Sprintf("WALLET_SEED=%s", info.Seed))
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_NAME=%s", info.Name))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ENDPOINT=%s", info.Endpoint))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ADMIN_PORT=%s", info.AdminPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_PORT=%s", info.AgentPort))

	//Janus-controller
	cmd.Env = append(cmd.Env, fmt.Sprintf("CONTROLLER_PORT=%s", info.ControllerPort))
}

func InstantiateAgent(agent ComposeInfo, profile string) error {
	command := "docker "

	// add -H hostname if remote deploying
	if agent.HostName != "" {
		command += fmt.Sprintf("-H ssh://%s ", agent.HostName)
	}

	// append the rest of the command
	command += fmt.Sprintf("compose -f /tmp/janus/docker-compose.yml --profile %s -p janus-%s up -d --no-recreate", profile, profile)
	logger.InfoLogger("Executing the %s command", command)

	//parse the command to a []string format and pass it to a CMD
	parsedCommand := helper.ParseCommand(command)
	cmd := exec.Command(parsedCommand[0], parsedCommand[1:]...)

	//parse compose info to environment variables
	parseComposeInfo(agent, cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
