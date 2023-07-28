package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	//"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
)

var (
	port = flag.Int("port", 12346, "The server port")
)

type PlantStorage struct {
	pb.UnimplementedPlantStorageServer
	db *gorm.DB
}

const (
	secretToken = "secert_token"
)

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
	errMigration := db.AutoMigrate(&models.Plant{}, &models.HumidityEntry{}, &models.ImageEntry{}, &models.Gpio{})
	if errMigration != nil {
		log.Fatalf("Unable to perform database migration. Terminating with error: %v", err)
		return
	}
	rpcServer(db)
}

func rpcServer(db *gorm.DB) {
	flag.Parse()
	localIp := "0.0.0.0" //GetOutboundIP().String()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(tokenInterceptor))
	pb.RegisterPlantStorageServer(s, &PlantStorage{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func tokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
	}

	tokenString := authHeader[0]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretToken), nil
	})
	if err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return handler(ctx, req)
}

func (s *PlantStorage) AddNewPlant(_ context.Context, request *pb.AddPlantRequest) (*pb.PlantOverviewMsg, error) {
	var plant models.Plant
	result := s.db.Model(&models.Plant{}).
		Where(models.Plant{Plant: request.PlantId}).
		First(&plant)

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

	for _, v := range request.Images {
		path := fmt.Sprintf("%s%d", "plant_", plant.Plant)
		resPath := utils.StoreImageInNewFile(v.ImageBytes, path, v.ImageId, true)
		imageEntry := models.ImageEntry{Plant: plant.Plant, Path: resPath}
		errCreateImage := s.db.Model(&models.ImageEntry{}).Create(&imageEntry).Error
		if errCreateImage != nil {
			log.Fatalf("Error: Unable to create the new Image. Errormessage: %s", errCreateImage.Error())
		} else {
			log.Println("New image added for plant: ", request.PlantId)
		}
	}

	// todo return thumbnail
	return &pb.PlantOverviewMsg{
		PlantId:   plant.Plant,
		Name:      plant.Name,
		Info:      plant.Info,
		Gpio:      &pb.GpioIdentifierMsg{GpioName: "tba"},
		Thumbnail: nil}, nil
}

func (s *PlantStorage) DeletePlant(_ context.Context, request *pb.PlantRequest) (*pb.DeletePlantReply, error) {
	errGetPlant := s.db.Model(&models.Plant{}).Delete(&models.Plant{}, request.Plant).Error
	if errGetPlant != nil {
		log.Fatalf("Error: Plant with Id: %d could not be deleted. Errormessage: %s", request.Plant, errGetPlant.Error())
	}
	return &pb.DeletePlantReply{}, nil
}

func (s *PlantStorage) GetOverviewAllPlants(_ context.Context, _ *pb.GetAllPlantsRequest) (*pb.AllPlantsReply, error) {
	var plants []models.Plant
	result := s.db.Find(&plants)
	if result.Error != nil {
		log.Fatalf("Error: Not able to query all plants. Errormessage: %s", result.Error.Error())
	}

	convertedPlants := make([]*pb.PlantOverviewMsg, len(plants))

	for i, v := range plants {
		convertedPlants[i] = &pb.PlantOverviewMsg{
			PlantId:   v.Plant,
			Name:      v.Name,
			Info:      v.Info,
			Gpio:      nil,
			Thumbnail: utils.LoadImageBytesFromPath(fmt.Sprintf("Storage/plants/plant_%d/0_thumbnail.jpg", v.Plant)),
		}
	}
	return &pb.AllPlantsReply{Plants: convertedPlants}, nil
}

func (s *PlantStorage) GetAdditionalDataPlant(_ context.Context, request *pb.GetAdditionalDataPlantRequest) (*pb.GetAdditionalDataPlantReply, error) {
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
	for i, v := range humidityEntries {
		convertedHumidity[i] = &pb.HumidityMsg{
			Humidity:  v.HumidityEntry,
			Timestamp: timestamppb.New(v.Timestamp),
		}
	}

	var imageEntries []models.ImageEntry
	errImage := s.db.Where(models.ImageEntry{Plant: request.PlantId}).Find(&imageEntries).Error
	if errImage != nil {
		log.Fatalf("Error: Plant with Id: %d unable to query Images entries. Errormessage: %s", request.PlantId, err.Error())
	}
	convertedImages := make([]*pb.ImageMsg, len(imageEntries))
	for i, v := range imageEntries {
		convertedImages[i] = &pb.ImageMsg{
			ImageId:    v.ImageEntry,
			ImageBytes: utils.LoadImageBytesFromPath(v.Path),
		}
	}

	return &pb.GetAdditionalDataPlantReply{Plant: request.PlantId, Humidity: convertedHumidity, Images: convertedImages}, nil
}

func (s *PlantStorage) StoreHumidityEntry(_ context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
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
