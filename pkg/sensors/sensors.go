package sensors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func CollectSensorData(ip, port string) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s", ip, port))
	if err != nil {
		fmt.Println(err)
	}

	data := make(map[string]float32)
	//{"temperature": 10, "humidity": 68} -> type
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)

	names := reflect.ValueOf(data).MapKeys()
	fmt.Println(names)
}
