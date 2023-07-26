package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	//"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"time"
)

var (
	port = flag.Int("port", 12346, "The server port")
)

type PlantStorage struct {
	pb.UnimplementedPlantStorageServer
	db *gorm.DB
}

func main() {
	//s := "user:password@tcp(0.0.0.0:3306)/mydatabase"
	//	db, err := gorm.Open(mysql.Open(s), &gorm.Config{})
	//s := ""
	s := "test.db"
	db, err := gorm.Open(sqlite.Open(s), &gorm.Config{})
	if err != nil {
		log.Fatalf("Terminating with error: %v", err)
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Plant{}, &models.HumidityEntry{})

	rpcServer(db)

}

func rpcServer(db *gorm.DB) {

	flag.Parse()
	localIp := "0.0.0.0" //GetOutboundIP().String()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPlantStorageServer(s, &PlantStorage{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *PlantStorage) StoreHumidityEntry(ctx context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	var plant *models.Plant
	errGetPlant := s.db.Model(&models.Plant{}).Where("plant = ?", request.RequestNumber).First(&plant).Error
	if errGetPlant != nil {
		log.Fatalf("Error: Plant with Id: %d could not be queried. Does it exist? Errormessage: %s", request.RequestNumber, errGetPlant.Error())
	}

	err := s.db.Model(&models.HumidityEntry{}).Create(&models.HumidityEntry{
		Plant:     plant.Plant,
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

func (s *PlantStorage) AddNewPlant(ctx context.Context, request *pb.AddPlantRequest) (*pb.PlantOverviewMsg, error) {
	var plant models.Plant
	result := s.db.Model(&models.Plant{}).
		First(&plant).
		Where(models.Plant{Plant: request.PlantId})
	if result.Error != nil {
		log.WithError(result.Error).Error("Error while querying existing plants.")
	}
	var errCreatePlant error
	if result.RowsAffected > 0 {
		log.Println("Existing plant will be updated")
		plant.Name = request.Name
		plant.Info = request.Info
		plant.SensorSlot = request.GpioSensorSlot
		s.db.Save(&plant)
	} else {
		plant := models.Plant{
			Name:       request.Name,
			Info:       request.Info,
			SensorSlot: request.GpioSensorSlot,
		}
		errCreatePlant = s.db.Model(&models.Plant{}).Create(&plant).Error
		log.Println("New plant added")
	}

	if errCreatePlant != nil {
		log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", errCreatePlant.Error())
	}
	// todo return thumbnail
	return &pb.PlantOverviewMsg{
		PlantId:   plant.Plant,
		Name:      plant.Name,
		Info:      plant.Info,
		Gpio:      &pb.GpioIdentifierMsg{GpioName: "tba"},
		Thumbnail: nil}, nil
}

func (s *PlantStorage) DeletePlant(ctx context.Context, request *pb.PlantRequest) (*pb.DeletePlantReply, error) {
	errGetPlant := s.db.Model(&models.Plant{}).Delete(&models.Plant{}, request.Plant).Error
	if errGetPlant != nil {
		log.Fatalf("Error: Plant with Id: %d could not be deleted. Errormessage: %s", errGetPlant.Error())
	}
	return &pb.DeletePlantReply{}, nil
}

func (s *PlantStorage) GetOverviewAllPlants(ctx context.Context, request *pb.GetAllPlantsRequest) (*pb.AllPlantsReply, error) {
	var plants []models.Plant
	result := s.db.Find(&plants)
	if result.Error != nil {
		log.Fatalf("Error: Not able to query all plants. Errormessage: %s", result.Error.Error())
	}

	convertedPlants := make([]*pb.PlantOverviewMsg, len(plants))

	for i, v := range convertedPlants {
		convertedPlants[i] = &pb.PlantOverviewMsg{
			PlantId:   v.PlantId,
			Name:      v.Name,
			Info:      v.Info,
			Gpio:      nil,
			Thumbnail: nil,
		}
	}
	return &pb.AllPlantsReply{Plants: convertedPlants}, nil
}

func (s *PlantStorage) GetAdditionalDataPlant(ctx context.Context, request *pb.GetAdditionalDataPlantRequest) (*pb.GetAdditionalDataPlantReply, error) {
	var plant models.Plant
	err := s.db.Where(models.Plant{Plant: request.PlantId}).FirstOrInit(&plant).Error
	if err != nil {
		log.Fatalf("Error: Plant with Id: %d does not exist yet. Errormessage: %s", request.PlantId, err.Error())
	}

	var humidityEntries []models.HumidityEntry
	errHumidity := s.db.Where(models.HumidityEntry{Plant: request.PlantId}).Find(&humidityEntries).Error
	if errHumidity != nil {
		log.Fatalf("Error: Plant with Id: %d unable to query Humidity entries. Errormessage: %s", request.PlantId, err.Error())
	}
	convertedHumidity := make([]*pb.HumidityMsg, len(humidityEntries))
	for i, v := range convertedHumidity {
		convertedHumidity[i] = &pb.HumidityMsg{
			Humidity:  v.Humidity,
			Timestamp: v.Timestamp,
		}
	}

	return &pb.GetAdditionalDataPlantReply{Plant: request.PlantId, Humidity: convertedHumidity}, nil
}

func (s *PlantStorage) getRequestedSensorStates(ctx context.Context, request *pb.GetRequestedSensorStatesRequest) (*pb.GetRequestedSensorStatesResponse, error) {
	//request.DeviceId
	/*var plant models.Plant
	result := s.db.Model(&models.Plant{}).
		First(&plant).
		Where(models.Plant{Plant: request.PlantId})
	if result.Error != nil {
		log.WithError(result.Error).Error("Error while querying existing plants.")
	}*/
	return &pb.GetRequestedSensorStatesResponse{}, nil
}
