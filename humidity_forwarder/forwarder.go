package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/receive", receiveHandler)
	log.Info("Starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Unable to start Server: ", err)
		return
	}
}

type Message struct {
	ID    int32 `json:"deviceid"`
	Value int32 `json:"value"`
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Received connection")
	if r.Method == "POST" {
		var message Message
		errJson := json.NewDecoder(r.Body).Decode(&message)
		if errJson != nil {
			http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error("Error while closing the http body: ", err.Error())
			}
		}(r.Body)

		log.Info("Received ID:", message.ID)
		log.Info("Received Value:", message.Value)
		forwardData(message.ID, message.Value)
	} else {
		http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
	}
}

func forwardData(id int32, value int32) {
	ForwardToHA(id, value)
	ForwardToPlantServer(id, value)
}
