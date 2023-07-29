package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"time"
)

func (s *PlantStorage) StoreHumidityEntry(_ context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	var sensor models.Sensor
	result := s.db.Model(&models.Sensor{}).
		Where(models.Sensor{SensorSlot: request.SensorId, DeviceId: request.DeviceId}).
		First(&sensor).Error
	sensorId := sensor.Sensor
	if result != nil {
		log.WithError(result).Error("Creating new Sensor.")
		sensor := models.Sensor{
			SensorSlot: request.SensorId,
			DeviceId:   request.DeviceId,
			InUse:      false,
		}
		errCreateSensor := s.db.Model(&models.Sensor{}).Create(&sensor).Error
		if errCreateSensor != nil {
			log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", errCreateSensor.Error())
		}

		sensorId = sensor.Sensor
	}

	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		Sensor:    sensorId,
		Value:     request.Humidity,
		Timestamp: time.Now(),
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
