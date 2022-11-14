package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
	"server/models"
	"time"
)

var (
	port = flag.Int("port", 12346, "The server port")
)

type StorageServer struct {
	pb.UnimplementedStorageServerServer
	db *gorm.DB
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Plant{}, &models.HumidityEntry{})

	createPlant(db)
	rpcServer(db)

}

func createPlant(db *gorm.DB) {
	errCreatePlant := db.Model(&models.Plant{}).Create(&models.Plant{
		Name:       "Testplant",
		Humidity:   1,
		SensorSlot: 14,
	}).Error
	if errCreatePlant != nil {
		log.Fatalf("Error while creating a new dish name rating.")
	}
}

func rpcServer(db *gorm.DB) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStorageServerServer(s, &StorageServer{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *StorageServer) StoreHumidityEntry(ctx context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	var plant *models.Plant
	errGetPlant := s.db.Model(&models.Plant{}).Where("plant = ?", request.PlantId).First(&plant).Error
	if errGetPlant != nil {
		log.Fatalf("Error: Plant with Id: %d could not be queried. Does it exist? Errormessage: %s", request.PlantId, errGetPlant.Error())
	}

	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		Plant:     plant.Plant,
		Value:     request.Humidity,
		Timestamp: time.Now(),
	}).Error
	if err != nil {
		log.Fatalf("Error: New Humidity Entry for Plant %d with value %d was not created. Errormessage: %s", request.PlantId, request.GetHumidity(), err.Error())
	} else {
		log.Println("New Humidity Entry for Plant %i with value %i", request.PlantId, request.GetHumidity())
	}
	return &pb.StoreHumidityReply{}, nil
}
