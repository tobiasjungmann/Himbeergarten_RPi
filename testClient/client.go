package main

import (
	//pb "../server/proto"
	log "github.com/sirupsen/logrus"
	//pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
}

// https://stackoverflow.com/questions/61311436/how-to-fix-parsing-go-mod-module-declares-its-path-as-x-but-was-required-as-y
// go mod edit -replace="github.com/y/original=github.com/x/version@latest
// go mod edit -replace example.com/tobiasjungmann/Himbeergarten_RPi/server=../server/proto
