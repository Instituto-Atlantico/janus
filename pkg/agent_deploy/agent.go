package agent_deploy

import (
	"fmt"
	"os"
	"os/exec"
)

type AgentInfo struct {
	Seed      string
	Name      string
	Endpoint  string
	AdminPort string
	AgentPort string
}

func InstantiateAgent(agent AgentInfo, hostName string) error {
	command := "docker "

	// add -H host name if remote deploying
	if hostName != "" {
		command += fmt.Sprintf("-H ssh://%s ", hostName)
	}

	// append the rest of the command
	command += "compose -f /tmp/janus/docker-compose.yml -p janus-agent up -d"

	parsedCommand := parseCommand(command)

	cmd := exec.Command(parsedCommand[0], parsedCommand[1:]...)

	cmd.Env = os.Environ()

	// The arguments are passed to the docker-compose.yml as env variables
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_PORT=%s", agent.AgentPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ADMIN_PORT=%s", agent.AdminPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("WALLET_SEED=%s", agent.Seed))
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_NAME=%s", agent.Name))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ENDPOINT=%s", agent.Endpoint))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
