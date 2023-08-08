package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

type HumidityData struct {
	State      int        `json:"state"`
	Attributes Attributes `json:"attributes"`
}

type HumidityDataInPercent struct {
	State      string     `json:"state"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	UnitOfMeasurement string `json:"unit_of_measurement"`
	FriendlyName      string `json:"friendly_name"`
}

func ForwardToHA(deviceId string, sensorId int32, humidity int32, humidityInPercent int32) {
	macForURL := strings.Replace(deviceId, ":", "_", 5)
	url := fmt.Sprintf("http://%s:8123/api/states/sensor.esp_%s_hum_%d_", *ipHa, macForURL, sensorId)
	log.Info("HA Address: ", url)
	sendRequest(fmt.Sprintf("%sraw", url), generatePayloadHumidity(humidity))
	sendRequest(fmt.Sprintf("%spercent", url), generatePayloadHumidityPercent(humidityInPercent))

}

func sendRequest(url string, humidity []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(humidity))

	if err != nil {
		log.Error("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("HOME_ASSISTANT_TOKEN"))
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

func generatePayloadHumidity(humidity int32) []byte {
	humidityData := HumidityData{
		State: int(humidity),
		Attributes: Attributes{
			FriendlyName: "Raw Humidity Data",
		},
	}

	jsonData, err := json.Marshal(humidityData)
	if err != nil {
		log.Error("Error marshaling the struct to JSON:", err)
		return []byte("")
	}
	return jsonData
}

func generatePayloadHumidityPercent(humidity int32) []byte {
	humidityData := HumidityDataInPercent{
		State: fmt.Sprint(humidity),
		Attributes: Attributes{
			UnitOfMeasurement: "%",
			FriendlyName:      "Humidity Data In Percent",
		},
	}

	jsonData, err := json.Marshal(humidityData)
	if err != nil {
		log.Error("Error marshaling the struct to JSON:", err)
		return []byte("")
	}
	return jsonData
}
