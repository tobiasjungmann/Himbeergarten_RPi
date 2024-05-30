package main

import log "github.com/sirupsen/logrus"

const (
	testPlantStorage  = false
	testForwarderGRPC = false
	testForwarderWeb  = true
)

func main() {
	if testPlantStorage {
		log.Info("Plant Management Tests Started")
		addEditPlantsOnPlantStorage()
	}
	if testForwarderGRPC {
		log.Info("Forwarder GRPC Tests Started")
		sendHumidityToForwarderGRPC()
	}
	if testForwarderWeb {
		log.Info("Forwarder Web Requests Tests Started")
		sendHumidityToForwarderWebRequest()
	}
}
