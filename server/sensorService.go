package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
	"net"
	"time"
)

type server struct {
	pb.UnsafeHumidityStorageServer
	db *gorm.DB
}

func startSensorAPI(db *gorm.DB) {
	localIp := "0.0.0.0"
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, 12347))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if len(*sslCertPath) > 0 && len(*sslKeyPath) > 0 {
		creds, err := credentials.NewServerTLSFromFile(*sslCertPath, *sslKeyPath)
		if err != nil {
			log.Fatalf("failed to load TLS certificates: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
		log.Info("TLS is activated.")
	} else {
		opts = []grpc.ServerOption{}
		log.Info("TLS is deactivated.")
	}

	s := grpc.NewServer(opts...)
	pb.RegisterHumidityStorageServer(s, &server{db: db})
	log.Info("Humidity Server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s server) GetActiveSensorsForDevice(_ context.Context, request *pb.GetActiveSensorsRequest) (*pb.GetActiveSensorsReply, error) {
	var sensor models.Sensor
	result := s.db.Model(&models.Sensor{}).
		Where(models.Sensor{SensorSlot: *request.SensorId, DeviceMAC: *request.DeviceMAC}).
		First(&sensor).Error
	sensorId := sensor.Sensor
	if result != nil {
		log.WithError(result).Error("Creating new Sensor.")
		sensor := models.Sensor{
			SensorSlot: *request.SensorId,
			DeviceMAC:  *request.DeviceMAC,
			InUse:      false,
		}
		errCreateSensor := s.db.Model(&models.Sensor{}).Create(&sensor).Error
		if errCreateSensor != nil {
			log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", errCreateSensor.Error())
		}

		sensorId = sensor.Sensor
	}

	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		Sensor:         sensorId,
		Value:          *request.Humidity,
		ValueInPercent: *request.HumidityInPercent,
		Timestamp:      time.Now(),
	}).Error
	if err != nil {
		log.Fatalf("Error: New Humidity Entry for Plant %d with value %d was not created. Errormessage: %s", sensorId, request.GetHumidity(), err.Error())
	} else {
		log.Println("New Humidity Entry for Plant %i with value %i", sensorId, request.GetHumidity())
	}
	return &pb.StoreHumidityReply{}, nil
}

func (s server) StoreHumidityEntry(_ context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	var sensor models.Sensor
	result := s.db.Model(&models.Sensor{}).
		Where(models.Sensor{SensorSlot: *request.SensorId, DeviceMAC: *request.DeviceMAC}).
		First(&sensor).Error
	sensorId := sensor.Sensor
	if result != nil {
		log.WithError(result).Error("Creating new Sensor.")
		sensor := models.Sensor{
			SensorSlot: *request.SensorId,
			DeviceMAC:  *request.DeviceMAC,
			InUse:      false,
		}
		errCreateSensor := s.db.Model(&models.Sensor{}).Create(&sensor).Error
		if errCreateSensor != nil {
			log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", errCreateSensor.Error())
		}

		sensorId = sensor.Sensor
	}

	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		Sensor:         sensorId,
		Value:          *request.Humidity,
		ValueInPercent: *request.HumidityInPercent,
		Timestamp:      time.Now(),
	}).Error
	if err != nil {
		log.Fatalf("Error: New Humidity Entry for Plant %d with value %d was not created. Errormessage: %s", sensorId, request.GetHumidity(), err.Error())
	} else {
		log.Println("New Humidity Entry for Plant %i with value %i", sensorId, request.GetHumidity())
	}
	return &pb.StoreHumidityReply{}, nil
}
