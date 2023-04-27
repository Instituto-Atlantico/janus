package sensors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CollectSensorData(ip, port string) map[string]any {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s", ip, port))
	if err != nil {
		log.Println(err)
	}

	//{"temperature": "10", "humidity": "68"}

	data := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	return data
}
