package main

import (
	//pb "../server/proto"
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	localAddress = "127.0.0.1:12346"
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
	defer cancel()
}

func createPlant(c pb.PlantStorageClient, ctx context.Context) {
	res, err := c.AddNewPlant(ctx, &pb.AddPlantRequest{
		PlantId:        0,
		Name:           "Test Pflanze 1",
		Info:           "Test Info 1",
		GpioSensorSlot: 0,
		Images:         nil,
	})

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Create plant request successful.")
		log.Info(res)
	}
}
