package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
)

func handleProtoBasedSocket() {
	localIp := "0.0.0.0"
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *apiPort))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			log.Info("Error while closing the socket: ", err.Error())
		}
	}(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			err := conn.Close()
			if err != nil {
				log.Info("Error while closing the connection: ", err.Error())
			}
			continue
		}
		log.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	/*defer func() {
		log.Info("closing connection")
		if err := conn.Close(); err != nil {
			log.Info("error closing connection:", err)
		}
	}()*/

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Info(err)
		return
	}
	if n <= 0 {
		log.Info("no data received")
		return
	}

	var storeHumidity pb.StoreHumidityRequest
	if err := proto.Unmarshal(buf[:n], &storeHumidity); err == nil {
		//log.Info("failed to unmarshal:", err)
		forwardData(storeHumidity.GetDeviceMAC(), storeHumidity.GetSensorId(),
			storeHumidity.GetHumidity(), storeHumidity.GetHumidityInPercent())
		return
	} else {
		var activeSensors pb.GetActiveSensorsRequest
		if err := proto.Unmarshal(buf[:n], &activeSensors); err == nil {
			log.Info("Received Active Sensor request from ", *activeSensors.DeviceMAC)

			msgInBytes, e := proto.Marshal(getActiveSensorsForDevice(&activeSensors))
			if e != nil {
				log.Info("Error while answering the client: ", err)
				return
			}
			_, err := conn.Write(msgInBytes)
			if err != nil {
				log.Info("Error while answering the client: ", err)
				return
			}
			return
		}
	}
	log.Info("failed to unmarshal - Message did not match any Message Format:", err)
}

func getActiveSensorsForDevice(sensors *pb.GetActiveSensorsRequest) *pb.GetActiveSensorsReply {
	address := fmt.Sprintf("%s:%d", *ipStorage, portStorage)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Info("Forwarder connecting to ", address)
	if err != nil {
		log.Error(err)
	}
	c := pb.NewHumidityStorageClient(conn)
	s, _ := generateToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", s)

	res, errStore := c.GetActiveSensorsForDevice(ctx, sensors)

	if errStore != nil {
		log.Error(errStore.Error())
	}
	return res
}
