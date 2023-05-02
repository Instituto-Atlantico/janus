package service

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
