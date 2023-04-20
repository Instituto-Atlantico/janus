package service

type ProvisionBody struct {
	DeviceHostName string   `json:"deviceHostName"` //user@ip
	DojotId        string   `json:"dojotId"`        //abc123
	Permissions    []string `json:"permissions"`    // ["temperature", "humidity"]
}

func ProvisionBodyIsValid(body ProvisionBody) bool {
	if body.DeviceHostName == "" || body.DojotId == "" || len(body.Permissions) == 0 {
		return false
	}

	return true
}
