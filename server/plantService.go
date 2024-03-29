package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
)

func (s *PlantStorage) AddNewPlant(_ context.Context, request *pb.AddPlantRequest) (*pb.PlantOverviewMsg, error) {
	var plant models.Plant
	result := s.db.Model(&models.Plant{}).
		Where(models.Plant{Plant: request.PlantId}).
		First(&plant)

	if result.Error != nil {
		log.WithError(result.Error).Error("Error while querying existing plants.")
	}

	// todo check that the sensor does actually exist
	var err error
	if result.RowsAffected > 0 {
		log.Println("Existing plant will be updated")
		plant.Name = request.Name
		plant.Info = request.Info
		plant.Sensor = request.SensorId
		s.db.Save(&plant)
	} else {
		plant := models.Plant{
			Name:   request.Name,
			Info:   request.Info,
			Sensor: request.SensorId,
		}
		err = s.db.Model(&models.Plant{}).Create(&plant).Error
		log.Println("New plant added")
	}

	if err != nil {
		log.Fatalf("Error: Unable to create the new Plant. Errormessage: %s", err.Error())
	}

	for _, v := range request.Images {
		path := fmt.Sprintf("%s%d", "plant_", plant.Plant)
		resPath := utils.StoreImageInNewFile(v.ImageBytes, path, v.ImageId, true)
		imageEntry := models.ImageEntry{Plant: plant.Plant, Path: resPath}

		if err := s.db.Model(&models.ImageEntry{}).Create(&imageEntry).Error; err != nil {
			log.Fatalf("Error: Unable to create the new Image. Errormessage: %s", err.Error())
		} else {
			log.Println("New image added for plant: ", request.PlantId)
		}
	}

	return &pb.PlantOverviewMsg{
		PlantId:   plant.Plant,
		Name:      plant.Name,
		Info:      plant.Info,
		Gpio:      &pb.GpioIdentifierMsg{GpioName: "tba"},
		Thumbnail: utils.LoadImageBytesFromPath(fmt.Sprintf("Storage/plants/plant_%d/0_thumbnail.jpg", plant.Plant)),
	}, nil
}

func (s *PlantStorage) DeletePlant(_ context.Context, request *pb.PlantRequest) (*pb.DeletePlantReply, error) {
	if err := s.db.Delete(&models.Plant{}, request.Plant).Error; err != nil {
		log.Fatalf("Error: Plant with Id: %d could not be deleted. Errormessage: %s", request.Plant, err.Error())
	}
	if err := os.RemoveAll(fmt.Sprintf("Storage/plants/plant_%d/", request.Plant)); err != nil {
		log.Error("Unable to delete existing images for deleted plant %d: %s", request.Plant, err)
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

	if err := s.db.Where(models.Plant{Plant: request.PlantId}).FirstOrInit(&plant).Error; err != nil {
		log.Fatalf("Error: Plant with Id: %d does not exist yet. Errormessage: %s", request.PlantId, err.Error())
	}

	var humidityEntries []models.HumidityEntry
	if err := s.db.Where(models.HumidityEntry{Sensor: plant.Sensor}).Find(&humidityEntries).Error; err != nil {
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

	if err := s.db.Where(models.ImageEntry{Plant: request.PlantId}).Find(&imageEntries).Error; err != nil {
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
