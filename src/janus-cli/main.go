package main

import (
	_ "embed"
	"errors"
	"log"
	"os"

	"github.com/Instituto-Atlantico/janus/src/janus-cli/cmd"
)

// This is copying the docker-compose file to the cli directory
//go:generate cp ../../docker/docker-compose.yml ./

// This is embeding the docker-compose file to the binary code
//
//go:embed docker-compose.yml
var dockercompose string

func generateTempFiles() error {
	path := "/tmp/janus/"

	// Generate janus path on tmp
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) { // Check if path already existes
		err := os.Mkdir(path, os.ModePerm) // Create path
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Store docker-compose file in tmp/janus.
	// Dockercompose variable is a string with the full body of our docker-compose.yml
	err = os.WriteFile("/tmp/janus/docker-compose.yml", []byte(dockercompose), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := generateTempFiles()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
