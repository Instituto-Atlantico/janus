package dojot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func PublishMessage(brokerServerUrl, brokerUsername, brokerPassword, sensorDataApiUrl, publicationTopic string) {
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(brokerServerUrl + ":1883").SetClientID("dojot")

	// Set username
	opts.SetUsername(brokerUsername)
	// Set password
	opts.SetPassword(brokerPassword)
	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	//opts.HTTPHeaders.Add()

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	resp, err := http.Get(sensorDataApiUrl + ":5000/")
	if err != nil {
		log.Fatal(err)
	}

	sensorData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Publish a message
	token := client.Publish(publicationTopic, 1, false, sensorData)
	token.Wait()

	time.Sleep(6 * time.Second)

	// Disconnect
	client.Disconnect(250)
	time.Sleep(1 * time.Second)
}
