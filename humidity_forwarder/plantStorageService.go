package main

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	localAddress            = "0.0.0.0:12347"
	secretTokenPlantStorage = "secret_token"
)

func generateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretTokenPlantStorage))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ForwardToPlantServer(id int32, value int32) {
	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err)
	}
	c := pb.NewHumidityStorageClient(conn)
	s, _ := generateToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", s)

	_, errStore := c.StoreHumidityEntry(ctx, &pb.StoreHumidityRequest{RequestNumber: id, Humidity: value})

	if errStore != nil {
		log.Error(errStore.Error())
	}
}
