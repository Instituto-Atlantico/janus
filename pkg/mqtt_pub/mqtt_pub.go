package mqtt_pub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BrokerCredentials struct {
	Ip       string
	Username string
	Password string
}

func PublishMessage(credentials BrokerCredentials, sensorData map[string]any) error {
	publicationTopic := fmt.Sprintf("%s/attrs", credentials.Username)

	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("TOPIC: %s\n", msg.Topic())
		log.Printf("MSG: %s\n", msg.Payload())
	}

	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(credentials.Ip + ":1883").SetClientID("dojot")

	// Set username
	opts.SetUsername(credentials.Username)
	// Set password
	opts.SetPassword(credentials.Password)
	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	parsedSensorData, err := json.Marshal(sensorData)
	if err != nil {
		return err
	}

	// Publish a message
	token := client.Publish(publicationTopic, 1, false, parsedSensorData)

	// time.Sleep(6 * time.Second)
	<-token.Done()

	// Disconnect
	client.Disconnect(250)
	time.Sleep(1 * time.Second)

	if token.Error() != nil {
		return token.Error()
	}
	return nil
}
