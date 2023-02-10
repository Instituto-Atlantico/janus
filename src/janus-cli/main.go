package main

import (
	_ "embed"
	"errors"
	"log"
	"os"

	"github.com/Instituto-Atlantico/janus/src/janus-cli/cmd"
)

//this is copying
//go:generate cp ../../docker/docker-compose.yml ./

// this is embeding the docker-compose file to the binary code
//
//go:embed docker-compose.yml
var dockercompose string

func generateTempFiles() {
	path := "/tmp/janus/"

	//generate janus path on tmp
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) { //check if path already exists
		err := os.Mkdir(path, os.ModePerm) // create path
		if err != nil {
			log.Println(err)
		}
	} else if err != nil {
		log.Println(err)
	}

	//store docker-compose file in tmp/janus.
	//dockercompose variable is a string with the full body of our docker-compose.yml
	err = os.WriteFile("/tmp/janus/docker-compose.yml", []byte(dockercompose), 0644)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	generateTempFiles()
	cmd.Execute()
}
