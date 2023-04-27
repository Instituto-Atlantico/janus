package main

import (
	_ "embed"
	"flag"
	"log"

	"github.com/Instituto-Atlantico/janus/src/janus-controller/service"
)

var (
	serverAgentIp string
	brokerIp      string
	port          string
	collectorTime int
)

func handleFlags() {
	flag.StringVar(&serverAgentIp, "server-agent-ip", "", "")
	flag.StringVar(&brokerIp, "broker-ip", "localhost", "The ip the broker is running")
	flag.StringVar(&port, "port", "8080", "The port the controller api will run")
	flag.IntVar(&collectorTime, "collector-time", 30, "The time between sensor collections")
	flag.Parse()

	if serverAgentIp == "" {
		log.Fatal("Required flag --server-agent-ip not passed")
	}
}

func main() {
	handleFlags()

	service := service.Service{}
	service.Init(serverAgentIp, brokerIp)
	service.RunCollector(collectorTime)
	service.RunApi(port)
}
