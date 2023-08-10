package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

func handleGRPC() {
	localIp := "0.0.0.0"
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *apiPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if len(*sslCertPath) > 0 && len(*sslKeyPath) > 0 {
		creds, err := credentials.NewServerTLSFromFile(*sslCertPath, *sslKeyPath)
		if err != nil {
			log.Fatalf("failed to load TLS certificates: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
		log.Info("TLS is activated.")
	} else {
		opts = []grpc.ServerOption{}
		log.Info("TLS is deactivated.")
	}

	s := grpc.NewServer(opts...)
	pb.RegisterHumidityStorageServer(s, &server{})
	log.Info("Forwarder RPC server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s server) StoreHumidityEntry(_ context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	forwardData(*request.DeviceMAC, *request.SensorId, *request.Humidity, *request.HumidityInPercent)
	return &pb.StoreHumidityReply{}, nil
}
