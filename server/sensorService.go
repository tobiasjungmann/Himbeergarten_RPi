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

/**
 * Return the GPIO numbers to which the sensors are connected for a specific device.
 * Devices are identified by their MAC address
 */
func (s server) GetActiveSensorsForDevice(_ context.Context, request *pb.GetActiveSensorsRequest) (*pb.GetActiveSensorsReply, error) {
	var sensor models.Sensor
	// check if all sensors are already stored in the database. If this is not hte case, they must be saved one by one
	err := s.db.Model(&models.Sensor{}).
		Where(models.Sensor{DeviceMAC: *request.DeviceMAC}).
		First(&sensor).Error

	var sensors []int32

	if err != nil {
		log.WithError(err).Error("Creating new Sensors.")
		for _, sens := range request.AvailableSensors {
			createSensor(sens, request.GetDeviceMAC(), s.db)
		}
	} else {
		// Device already exists - which sensors are already loaded? todo query for all which are set to active
		err = s.db.Model(&models.Sensor{}).
			Where(models.Sensor{DeviceMAC: *request.DeviceMAC, InUse: true}).Select("SensorSlot").Find(sensors).Error
		if err := s.db.Create(&sensor).Error; err != nil {
			log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", err.Error())
		}
	}

	return &pb.GetActiveSensorsReply{Sensors: sensors}, nil
}

func createSensor(sensorId int32, deviceMac string, db *gorm.DB) models.Sensor {
	sensor := models.Sensor{
		SensorSlot: sensorId,
		DeviceMAC:  deviceMac,
		InUse:      false,
	}

	if err := db.Create(&sensor).Error; err != nil {
		log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", err.Error())
	}
	return sensor
}

func (s server) StoreHumidityEntry(_ context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	var sensor models.Sensor
	result := s.db.Model(&models.Sensor{}).
		Where(models.Sensor{SensorSlot: *request.SensorId, DeviceMAC: *request.DeviceMAC}).
		First(&sensor).Error

	if result != nil {
		sensor = createSensor(request.GetSensorId(), request.GetDeviceMAC(), s.db)
	}

	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		Sensor:         sensor.Sensor,
		Value:          *request.Humidity,
		ValueInPercent: *request.HumidityInPercent,
		Timestamp:      time.Now(),
	}).Error

	if err != nil {
		log.Fatalf("Error: New Humidity Entry for Plant %d with value %d was not created. Errormessage: %s", sensor.Sensor, request.GetHumidity(), err.Error())
	} else {
		log.Println("New Humidity Entry for Plant %i with value %i", sensor.Sensor, request.GetHumidity())
	}
	return &pb.StoreHumidityReply{}, nil
}
