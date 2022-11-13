package main

import (
	"context"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"gorm.io/gorm"
)

type StorageServer struct {
	pb.UnimplementedStorageServer
	db        *gorm.DB
	deviceBuf *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
}

func (s *server) StoreHumidityEntry(ctx context.Context, in *StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	return &pb.StoreHumidityReply{}, nil
}
