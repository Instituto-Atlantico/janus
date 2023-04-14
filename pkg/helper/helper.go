package helper

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"time"
)

func TryUntilNoError[R any](fn func() (R, error), timeoutInSeconds int) (R, error) {
	cResponse := make(chan R)
	cTimeout := make(chan string)

	go func() {
		time.Sleep(time.Second * time.Duration(timeoutInSeconds))
		cTimeout <- ""
	}()

	go func() {
		for {
			response, err := fn()
			if err == nil {
				cResponse <- response
				return
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case data := <-cResponse:
		return data, nil
	case <-cTimeout:
		return *new(R), errors.New("Timeout")
	}
}

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
func ParseCommand(command string) []string {
	parsed := strings.Split(command, " ")
	return parsed
}

func SliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
