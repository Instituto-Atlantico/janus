package main

import (
	"fmt"

	"github.com.br/janus/agent_deploy/helper"
)

func main() {
	// - Discover IP
	ip, err := helper.GetOutboundIP()
	if err != nil {
		fmt.Println(err)
	}

	endpoint := fmt.Sprintf("http://%s", ip.String())
	fmt.Printf("Actual IP: %s\n", endpoint)

	// - Provision Port // needs to be done manually because it needs to be open in firewall
	agentPort := "8000"
	adminPort := "8001"

	fmt.Printf("Admin port: %s\nAgent port: %s\n", adminPort, agentPort)

	// - Generate seed
	seed := helper.GenerateSeed()

	fmt.Printf("Seed generated: %s\n", seed)

	// - Register DID
	did, err := helper.RegisterDID(seed, "http://dev.bcovrin.vonx.io")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("DiD registered: %s\n", did)

	// - Instatiate Agent
	err = helper.InstatiateAgent(seed, "vitor", adminPort, agentPort, ip.String())
	if err != nil {
		fmt.Println(err)
	}
}
