package main

import (
	_ "embed"
	"flag"
	"log"

	"github.com/Instituto-Atlantico/janus/src/janus-controller/service"
)

var (
	serverAgentIp string
	port          string
	collectorTime int
)

func handleFlags() {
	flag.StringVar(&serverAgentIp, "server-agent-ip", "", "")
	flag.StringVar(&port, "port", "8081", "The port the controller api will run")
	flag.IntVar(&collectorTime, "collector-time", 30, "The time between sensor collections")
	flag.Parse()

	if serverAgentIp == "" {
		log.Fatal("Required flag --server-agent-ip not passed")
	}
}

func main() {
	handleFlags()

	service := service.Service{}
	service.Init(serverAgentIp)
	service.RunCollector(collectorTime)
	service.RunApi(port)
}
