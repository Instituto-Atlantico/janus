package main

import (
	_ "embed"
	"log"

	"github.com/Instituto-Atlantico/janus/pkg/temp_files"
	"github.com/Instituto-Atlantico/janus/src/janus-cli/cmd"
)

// This is copying the docker-compose file to the cli directory
//go:generate cp ../../docker/docker-compose.yml ./

// This is embedding the docker-compose file to the binary code
//
//go:embed docker-compose.yml
var dockercompose string

func main() {
	err := temp_files.GenerateTempFiles(dockercompose)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
