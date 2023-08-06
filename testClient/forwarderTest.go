package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net/http"
	"time"
)

const (
	serverURL   = "http://localhost:8080/receive"
	addressGRPC = "0.0.0.0:12346"
	secretToken = "secret_token"
)

func sendHumidityToForwarderWebRequest() {
	data := map[string]int32{
		"deviceid": 42, //"te:st:id:js:on",
		"value":    32,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshaling data to JSON: ", err)
		return
	}

	// Make an HTTP POST request to the server
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error sending POST request: ", err)
		return
	}
	defer resp.Body.Close()

	// Check the response from the server
	if resp.StatusCode == http.StatusOK {
		log.Info("Data sent successfully to the server.")
	} else {
		log.Warning("Failed to send data to the server. Status code:", resp.StatusCode)
	}
}

func sendHumidityToForwarderGRPC() {
	conn, err := grpc.Dial(addressGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
