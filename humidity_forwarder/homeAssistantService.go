package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	testTokenHomeAssistant = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI4MzczNDExNWFlNDc0ZGY4YjJiOGRlNWEzMDZkNTFkMCIsImlhdCI6MTY5MDY0NzM1MSwiZXhwIjoyMDA2MDA3MzUxfQ.RBXcYVaGhas-GPBt-04jE56TX1X50E7ypTJIKR-7zYQ"
)

type HumidityData struct {
	State      string     `json:"state"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	UnitOfMeasurement string `json:"unit_of_measurement"`
	FriendlyName      string `json:"friendly_name"`
}

func ForwardToHA(id string, value int32) {
	url := fmt.Sprintf("http://%s:8123/api/states/sensor.humidity%s", *ipHa, id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(generatePayload(value))))

	if err != nil {
		log.Error("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+testTokenHomeAssistant)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending HTTP request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Error closing connection to home assistant:", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		log.Info("Request successful!")
	} else {
		log.Info("Request failed with status code:", resp.StatusCode)
	}
}

func generatePayload(value int32) []byte {
	humidityData := HumidityData{
		State: fmt.Sprint(value),
		Attributes: Attributes{
			UnitOfMeasurement: "%",
			FriendlyName:      "Humidity data Input 1",
		},
	}
	jsonData, err := json.Marshal(humidityData)
	if err != nil {
		log.Error("Error marshaling the struct to JSON:", err)
		return []byte("")
	}
	return jsonData
}
