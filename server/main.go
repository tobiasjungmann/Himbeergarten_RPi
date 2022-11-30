package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
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

	rpcServer(db)

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func rpcServer(db *gorm.DB) {

	flag.Parse()
	localIp := GetOutboundIP().String()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *port))
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

func (s *StorageServer) AddNewPlant(ctx context.Context, request *pb.AddPlantRequest) (*pb.PlantOverviewMsg, error) {
	// todo store additional images
	log.Println("New plant added")
	plant := models.Plant{
		Name:       request.Name,
		Info:       request.Info,
		SensorSlot: request.GpioSensorSlot,
	}
	errCreatePlant := s.db.Model(&models.Plant{}).Create(&plant).Error
	if errCreatePlant != nil {
		log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", errCreatePlant.Error())
	}
	// todo return thumbnail
	return &pb.PlantOverviewMsg{
		PlantId:        plant.Plant,
		Name:           plant.Name,
		Info:           plant.Info,
		GpioSensorSlot: plant.SensorSlot,
		Thumbnail:      nil}, nil
}

func (s *StorageServer) DeletePlant(ctx context.Context, request *pb.DeletePlantRequest) (*pb.DeletePlantReply, error) {
	errGetPlant := s.db.Model(&models.Plant{}).Delete(&models.Plant{}, request.Plant).Error
	if errGetPlant != nil {
		log.Fatalf("Error: Plant with Id: %d could not be deleted. Errormessage: %s", errGetPlant.Error())
	}
	return &pb.DeletePlantReply{}, nil
}

func (s *StorageServer) GetOverviewAllPlants(ctx context.Context, request *pb.GetAllPlantsRequest) (*pb.AllPlantsReply, error) {
	/*	errGetPlant := s.db.Model(&models.Plant{}).Delete(&models.Plant{}, request.Plant).Error
		if errGetPlant != nil {
			log.Fatalf("Error: Plant with Id: %d could not be deleted. Errormessage: %s", errGetPlant.Error())
		}

		// todo query and parse plants to the reply
		&pb.PlantOverviewMsg{Plant: plant.Plant, Name: plant.Name,
			Info:      plant.Info,
			Type:      plant.Type,
			Thumbnail: nil}
	*/
	return &pb.AllPlantsReply{}, nil
}

func (s *StorageServer) GetAdditionalDataPlant(ctx context.Context, request *pb.GetAdditionalDataPlantRequest) (*pb.GetAdditionalDataPlantReply, error) {
	/*errGetPlant := s.db.Model(&models.Plant{}).Delete(&models.Plant{}, request.Plant).Error
	if errGetPlant != nil {
		log.Fatalf("Error: Plant with Id: %d could not be deleted. Errormessage: %s", errGetPlant.Error())
	}*/
	return &pb.GetAdditionalDataPlantReply{}, nil
}
