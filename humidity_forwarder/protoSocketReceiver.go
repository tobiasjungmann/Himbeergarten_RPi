package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
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
	defer func() {
		log.Info("closing connection")
		if err := conn.Close(); err != nil {
			log.Info("error closing connection:", err)
		}
	}()

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

	var e pb.StoreHumidityRequest
	if err := proto.Unmarshal(buf[:n], &e); err != nil {
		log.Info("failed to unmarshal:", err)
		return
	}

	/*log.Printf("{DeviceID:%s, Humidity:%d}\n",
	e.GetDeviceId(),
	e.GetHumidity())*/
	forwardData(e.GetDeviceId(), e.GetSensorId(),
		e.GetHumidity(), e.GetHumidityInPercent())

}
