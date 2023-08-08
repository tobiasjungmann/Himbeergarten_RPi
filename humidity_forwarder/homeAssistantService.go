package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

const (
	testTokenHomeAssistant = ""
)

type HumidityData struct {
	State      string     `json:"state"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	UnitOfMeasurement string `json:"unit_of_measurement"`
	FriendlyName      string `json:"friendly_name"`
	ValueInPercent    string `valueInPercent`
}

func ForwardToHA(deviceId string, sensorId int32, humidity int32, humidityInPercent int32) {
	macForURL := strings.Replace(deviceId, ":", "_", 5)
	url := fmt.Sprintf("http://%s:8123/api/states/sensor.humidity%s_sensor%d", *ipHa, macForURL, sensorId)
	log.Info("HA Address: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(generatePayload(humidity, humidityInPercent))))

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
		log.Info("Request to Home Assistant failed with status code:", resp.StatusCode)
	}
}

func generatePayload(humidity int32, humidityInPercent int32) []byte {
	humidityData := HumidityData{
		State: fmt.Sprint(humidity),
		Attributes: Attributes{
			UnitOfMeasurement: "%",
			FriendlyName:      "Humidity data Input 1",
			ValueInPercent:    "76%",
		},
	}
	jsonData, err := json.Marshal(humidityData)
	if err != nil {
		log.Error("Error marshaling the struct to JSON:", err)
		return []byte("")
	}
	return jsonData
}
