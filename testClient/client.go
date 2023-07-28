package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	localAddress = "0.0.0.0:12346"
	testImage    = "./images/IMG_20221218_135005.jpg"
)

func main() {
	log.Info("Connecting...")

	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Info(err)
	}
	c := pb.NewPlantStorageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	createPlant(c, ctx)
	getPlantOverview(c, ctx)
	defer cancel()
}

func createPlant(c pb.PlantStorageClient, ctx context.Context) {

	images := make([]*pb.ImageMsg, 1)
	images[0] = &pb.ImageMsg{ImageBytes: utils.LoadImageBytesFromPath(testImage), ImageId: 0}
	res, err := c.AddNewPlant(ctx, &pb.AddPlantRequest{
		PlantId:        0,
		Name:           "Test Pflanze 1",
		Info:           "Test Info 1",
		GpioSensorSlot: 0,
		Images:         images,
	})

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Create plant request successful.")
		log.Info(res)
	}
}

func getPlantOverview(c pb.PlantStorageClient, ctx context.Context) {
	res, err := c.GetOverviewAllPlants(ctx, &pb.GetAllPlantsRequest{})

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Create plant request successful.")
		//log.Info(res)
		log.Info("Plants:")
		for _, v := range res.Plants {
			log.Info("Id:   ", v.PlantId)
			log.Info("Name: ", v.Name)
			log.Info("Info: ", v.Info)
			utils.StoreImageInNewFile(v.Thumbnail, fmt.Sprintf("images/test/plant_%d/0_thumbnail.jpg", v.PlantId), 0, false)
		}
	}
}
