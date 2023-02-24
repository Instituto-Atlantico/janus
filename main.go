package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Parses a string command to a slice filled with every single part of it.
func parseCommand(command string) []string {
	parsed := strings.Split(command, " ")
	return parsed
}

func main() {
	command := "docker -H ssh://pi@raspberrypi.local compose -f /tmp/janus/docker-compose.yml up -d"

	parsedCommand := parseCommand(command)

	fmt.Println(command)

	cmd := exec.Command(parsedCommand[0], parsedCommand[1:]...)

	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_PORT=%s", "8001"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ADMIN_PORT=%s", "8002"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("WALLET_SEED=%s", "9b3bku6liM8aEsYKtBD76jenZUKIa7ZU"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("AGENT_NAME=%s", "test"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ENDPOINT=%s", "raspberrypi.local"))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
