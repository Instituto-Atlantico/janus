package agent_deploy

import (
	"net"
	"strings"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err

	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

// parseCommand parses a string command to a slice filled with every single part of it.
func parseCommand(command string) []string {
	parsed := strings.Split(command, " ")
	return parsed
}
