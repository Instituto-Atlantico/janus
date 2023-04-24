package service

type ProvisionBody struct {
	DeviceHostName string   `json:"deviceHostName"` // user@ip
	Permissions    []string `json:"permissions"`    // ["temperature", "humidity"]
	BrokerServerIp string   `json:"brokerServerIp"` // brokerIP
	BrokerUsername string   `json:"brokerUsername"` // brokerUser:deviceID
	BrokerPassword string   `json:"brokerPassword"` // brokerPass
}

func ProvisionBodyIsValid(body ProvisionBody) bool {
	if body.DeviceHostName == "" || len(body.Permissions) == 0 || body.BrokerServerIp == "" || body.BrokerUsername == "" || body.BrokerPassword == "" {
		return false
	}

	return true
}
