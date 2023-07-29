package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	serverURL = "http://localhost:8080/receive"
)

func mockESPHumidity() {
	data := map[string]int32{
		"deviceid": 42,
		"value":    32,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshaling data to JSON: ", err)
		return
	}

	// Make an HTTP POST request to the server
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error sending POST request: ", err)
		return
	}
	defer resp.Body.Close()

	// Check the response from the server
	if resp.StatusCode == http.StatusOK {
		log.Info("Data sent successfully to the server.")
	} else {
		log.Warning("Failed to send data to the server. Status code:", resp.StatusCode)
	}
}
