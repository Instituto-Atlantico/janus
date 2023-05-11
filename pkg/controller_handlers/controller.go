package controller_handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProvisionBody struct {
	DeviceHostName string   `json:"deviceHostName"` // user@ip
	Permissions    []string `json:"permissions"`    // ["temperature", "humidity"]
	BrokerIp       string   `json:"brokerIp"`       // 127.0.0.1
	BrokerUsername string   `json:"brokerUsername"` // brokerUser:deviceID
	BrokerPassword string   `json:"brokerPassword"` // brokerPass
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
