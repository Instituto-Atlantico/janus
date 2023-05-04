package sensors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
)

func CollectSensorData(ip, port string) (map[string]any, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s", ip, port))
	connectionRefused := errors.Is(err, syscall.ECONNREFUSED)
	if err != nil && connectionRefused {
		err := errors.New("Error requesting sensor data. The API doesn't seem to be working properly. Check its availability")

		return nil, err
	}

	if err != nil {
		return nil, err
	}

	//{"temperature": "10", "humidity": "68"}
	data := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err := errors.New("Error in reading sensor data. Check its physical integrity")

		return nil, err
	}

	return data, nil
}
