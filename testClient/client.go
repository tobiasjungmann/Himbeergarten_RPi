package main

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	localAddress = "0.0.0.0:12346"
	testImage    = "./images/IMG_20221218_135005.jpg"
	secretToken  = "secret_token"
)

func generateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	log.Info("Connecting...")
	//testPlantStorage()
	testHumidityForwarder()
	//mockESPHumidity()
}

func testHumidityForwarder() {
	conn, err := grpc.Dial("192.168.0.6:12348", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Info(err)
	}
	c := pb.NewHumidityStorageClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background())
	var humidity int32 = 698
	var humidityInPercent int32 = 42
	var sensorId int32 = 0
	var deviceMAC = "te:st:MAC:00:00:00"
	res, err := c.StoreHumidityEntry(ctx, &pb.StoreHumidityRequest{
		Humidity:          &humidity,
		SensorId:          &sensorId,
		HumidityInPercent: &humidityInPercent,
		DeviceId:          &deviceMAC,
	})

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Request send successfully request successful.")
		log.Info(res)
	}

}

func testPlantStorage() {
	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Info(err)
	}
	c := pb.NewPlantStorageClient(conn)
	s, _ := generateToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", s)

	createPlant(c, ctx)
	getPlantOverview(c, ctx)
}

func createPlant(c pb.PlantStorageClient, ctx context.Context) {
	images := make([]*pb.ImageMsg, 1)
	images[0] = &pb.ImageMsg{ImageBytes: utils.LoadImageBytesFromPath(testImage), ImageId: 0}
	res, err := c.AddNewPlant(ctx, &pb.AddPlantRequest{
		PlantId: 0,
		Name:    "Test Pflanze 1",
		Info:    "Test Info 1",
		Images:  images,
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
