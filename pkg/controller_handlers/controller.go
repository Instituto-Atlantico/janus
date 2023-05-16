package controller_handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProvisionBody struct {
	DeviceHostName string   `json:"deviceHostName" example:"rasp@192.168.0.5"`   // user@ip
	Permissions    []string `json:"permissions" example:"temperature, humidity"` // ["temperature", "humidity"]
	BrokerIp       string   `json:"brokerIp" example:"192.168.0.2"`              // 127.0.0.1
	BrokerUsername string   `json:"brokerUsername" example:"admin:a1998e"`       // brokerUser:deviceID
	BrokerPassword string   `json:"brokerPassword" example:"admin"`              // brokerPass
}

func ProvisionBodyIsValid(body ProvisionBody) bool {
	if body.DeviceHostName == "" || body.BrokerIp == "" || len(body.Permissions) == 0 || body.BrokerUsername == "" || body.BrokerPassword == "" {
		return false
	}

	return true
}

type ControllerService struct {
	Address string
}

func (c ControllerService) RequestProvision(body ProvisionBody) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := fmt.Sprintf("http://%s/provision", c.Address)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("provision request failed with status %v", req.Response.StatusCode)
	}

	return nil
}
