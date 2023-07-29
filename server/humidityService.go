package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"time"
)

func (s *PlantStorage) StoreHumidityEntry(_ context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		DeviceID:   request.DeviceId,
		SensorSlot: request.SensorId,
		Value:      request.Humidity,
		Timestamp:  time.Now(),
	}).Error
	if err != nil {
		log.Fatalf("Error: New Humidity Entry for Plant %d with value %d was not created. Errormessage: %s", request.RequestNumber, request.GetHumidity(), err.Error())
	} else {
		log.Println("New Humidity Entry for Plant %i with value %i", request.RequestNumber, request.GetHumidity())
	}
	return &pb.StoreHumidityReply{}, nil
}

func (s *PlantStorage) getRequestedSensorStates(_ context.Context, _ *pb.GetRequestedSensorStatesRequest) (*pb.GetRequestedSensorStatesResponse, error) {
	/*	request.DeviceId
		var plant models.Plant
		result := s.db.Model(&models.Device{}).
			First(&plant).
			Where(models.Plant{Plant: request.PlantId})
		if result.Error != nil {
			log.WithError(result.Error).Error("Error while querying existing plants.")
		}*/
	return &pb.GetRequestedSensorStatesResponse{}, nil
}
