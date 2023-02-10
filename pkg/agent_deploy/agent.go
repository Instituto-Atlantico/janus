package agent_deploy

import (
	"fmt"
	"os"
	"os/exec"
)

func InstatiateAgent(seed, name, adminPort, agentPort, endpoint string) error {
	command := "docker compose -f /tmp/janus/docker-compose.yml -p janus-agent up -d"
	parsedCommand := parseCommand(command)

	cmd := exec.Command(parsedCommand[0], parsedCommand[1:]...)

	//the arguments are passed to the docker-compose.yml as env variables
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_PORT=%s", agentPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ADMIN_PORT=%s", adminPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("WALLET_SEED=%s", seed))
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_NAME=%s", name))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ENDPOINT=%s", endpoint))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
