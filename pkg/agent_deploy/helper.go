package agent_deploy

import (
	"net"
	"regexp"
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

// ValidateSSHHostName uses regex to validate the format user@hostname for ssh connections
func ValidateSSHHostName(hostName string) bool {
	re := regexp.MustCompile("(?i)[A-Za-z]+@[A-Za-z-z0-9]+")
	return re.MatchString(hostName)
}

// Parses a string command to a slice filled with every single part of it.
func parseCommand(command string) []string {
	parsed := strings.Split(command, " ")
	return parsed
}
