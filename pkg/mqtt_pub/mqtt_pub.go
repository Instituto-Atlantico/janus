package mqtt_pub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BrokerData struct {
	BrokerServerIp string
	BrokerPassword string
}

func PublishMessage(brokerData BrokerData, brokerUsername, publicationTopic string, sensorData map[string]any) {
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(brokerData.BrokerServerIp + ":1883").SetClientID("dojot")

	// Set username
	opts.SetUsername(brokerUsername)
	// Set password
	opts.SetPassword(brokerData.BrokerPassword)
	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	//opts.HTTPHeaders.Add()

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	parsedSensorData, err := json.Marshal(sensorData)
	if err != nil {
		log.Println(err)
	}

	// Publish a message
	token := client.Publish(publicationTopic, 1, false, parsedSensorData)
	token.Wait()

	time.Sleep(6 * time.Second)

	// Disconnect
	client.Disconnect(250)
	time.Sleep(1 * time.Second)
}
