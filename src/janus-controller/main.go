package main

import (
	_ "embed"
	"flag"
	"log"

	"github.com/Instituto-Atlantico/janus/pkg/temp_files"
	"github.com/Instituto-Atlantico/janus/src/janus-controller/service"
)

// This is copying the docker-compose file to the cli directory
//go:generate cp ../../docker/docker-compose.yml ./

// This is embedding the docker-compose file to the binary code
//
//go:embed docker-compose.yml
var dockercompose string

var (
	serverAgentIp string
	port          string
	collectorTime int
)

func handleFlags() {
	flag.StringVar(&serverAgentIp, "server-agent-ip", "", "")
	flag.StringVar(&port, "port", "8080", "")
	flag.IntVar(&collectorTime, "collector-time", 30, "")
	flag.Parse()

	if serverAgentIp == "" {
		log.Fatal("Required flag --server-agent-ip not passed")
	}
}

func main() {
	err := temp_files.GenerateTempFiles(dockercompose)
	if err != nil {
		log.Fatal(err)
	}

	handleFlags()

	service := service.Service{}
	service.Init(serverAgentIp)
	service.RunCollector(collectorTime)
	service.RunApi(port)
}
