package helper

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type registerDIDRequest struct {
	Alias string `json:"alias"`
	Seed  string `json:"seed"`
	Role  string `json:"role"`
	DID   string `json:"did"`
}

type RegisterDIDResponse struct {
	DID    string `json:"did"`
	Seed   string `json:"seed"`
	Verkey string `json:"verkey"`
}

func GenerateSeed() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length := 32

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RegisterDID(seed, ledgerURL string) (string, error) {
	request := registerDIDRequest{
		Alias: "",
		Seed:  seed,
		Role:  "",
	}

	var response RegisterDIDResponse

	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(ledgerURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.DID, err
}
