package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
)

var (
	port = flag.Int("port", 12346, "The server port")
)

type StorageServer struct {
	pb.UnimplementedStorageServerServer
	db *gorm.DB
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStorageServerServer(s, &StorageServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *StorageServer) StoreHumidityEntry(ctx context.Context, in *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	return &pb.StoreHumidityReply{}, nil
}
