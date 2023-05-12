package service

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
