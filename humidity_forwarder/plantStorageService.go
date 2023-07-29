package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func ForwardToPlantServer(id int32, value int32) {
	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err)
	}
	c := pb.NewPlantStorageClient(conn)
	s, _ := generateToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", s)

	_, errStore := c.StoreHumidityEntry(ctx, &pb.StoreHumidityRequest{RequestNumber: id, Humidity: value})

	if errStore != nil {
		log.Error(errStore.Error())
	}
}
