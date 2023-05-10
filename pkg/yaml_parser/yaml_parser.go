package yaml_parser

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Agent struct {
	Hostname string   `yaml:"hostname"`
	Sensors  []string `yaml:"sensors"`
	Broker   struct {
		IP       string `yaml:"ip"`
		ID       string `yaml:"id"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"broker"`
}

type Body struct {
	Default Agent   `yaml:"default"`
	Agents  []Agent `yaml:"agents"`
}

func checkAgentRequiredFields(agent Agent) error {
	emptyFields := []string{}
	if agent.Hostname == "" {
		emptyFields = append(emptyFields, "hostname")
	}
	if len(agent.Sensors) == 0 {
		emptyFields = append(emptyFields, "sensors")
	}
	if agent.Broker.IP == "" {
		emptyFields = append(emptyFields, "broker.id")
	}
	if agent.Broker.ID == "" {
		emptyFields = append(emptyFields, "broker.ip")
	}
	if agent.Broker.Username == "" {
		emptyFields = append(emptyFields, "broker.username")
	}
	if agent.Broker.Password == "" {
		emptyFields = append(emptyFields, "broker.password")
	}

	parsed, err := yaml.Marshal(agent)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if len(emptyFields) > 0 {
		return fmt.Errorf("missing fields:%s in agent:\n%s", emptyFields, string(parsed))
	}

	return nil
}

func validateBody(body Body) error {
	for _, agent := range body.Agents {
		err := checkAgentRequiredFields(agent)
		if err != nil {
			return err
		}
	}

	return nil
}

func fillAgent(agent, defaultAgent Agent) Agent {
	if len(agent.Sensors) == 0 {
		agent.Sensors = defaultAgent.Sensors
	}
	if agent.Broker.IP == "" {
		agent.Broker.IP = defaultAgent.Broker.IP
	}
	if agent.Broker.Username == "" {
		agent.Broker.Username = defaultAgent.Broker.Username
	}
	if agent.Broker.Password == "" {
		agent.Broker.Password = defaultAgent.Broker.Password
	}
	return agent
}

func fillAgents(body *Body) {
	for i, agent := range body.Agents {
		body.Agents[i] = fillAgent(agent, body.Default)
	}
}

func ParseFile(filepath string) Body {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	body := Body{}

	err = yaml.Unmarshal(yamlFile, &body)
	if err != nil {
		log.Fatal(err)
	}

	fillAgents(&body)
	err = validateBody(body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}
