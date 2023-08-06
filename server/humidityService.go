package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PlantStorage) getConnectedSensorOverview(_ context.Context, _ *pb.GetSensorOverviewRequest) (*pb.GetSensorOverviewResponse, error) {
	var sensors []models.Sensor
	result := s.db.Find(&sensors)
	if result.Error != nil {
		log.Fatalf("Error: Not able to query all sensors. Errormessage: %s", result.Error.Error())
	}

	convertedSensors := make([]*pb.SensorMsg, len(sensors))

	for i, v := range sensors {
		var plant models.Plant
		err := s.db.Where(models.Plant{Sensor: v.Sensor}).First(&plant).Error
		var id int32 = -1
		if err == nil {
			id = plant.Plant
		}
		convertedSensors[i] = &pb.SensorMsg{
			DeviceMAC:      v.DeviceMAC,
			SensorId:       v.Sensor,
			SensorSlot:     v.SensorSlot,
			InUse:          v.InUse,
			ConnectedPlant: id,
		}
	}
	return &pb.GetSensorOverviewResponse{Sensors: convertedSensors}, nil
}

func (s *PlantStorage) GetDataForSensor(_ context.Context, request *pb.GetDataForSensorRequest) (*pb.GetDataForSensorReply, error) {
	var humidityEntries []models.HumidityEntry

	errHumidity := s.db.Where(models.HumidityEntry{Sensor: request.SensorId}).Find(&humidityEntries).Error
	if errHumidity != nil {
		log.Fatalf("Error: Plant with Id: %d unable to query Humidity entries. Errormessage: %s", request.SensorId, errHumidity.Error())
	}
	convertedHumidity := make([]*pb.HumidityMsg, len(humidityEntries))
	for i, v := range humidityEntries {
		convertedHumidity[i] = &pb.HumidityMsg{
			Humidity:  v.HumidityEntry,
			Timestamp: timestamppb.New(v.Timestamp),
		}
	}

	return &pb.GetDataForSensorReply{Data: convertedHumidity}, nil
}
