package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"net/http"
	"os"
)

var (
	port        = flag.Int("port", 12348, "The forwarder port")
	portStorage = flag.Int("portStorage", 12347, "The forwarder port")
	ipStorage   = flag.String("ipStorage", "0.0.0.0", "The forwarder port")
	sslCertPath = flag.String("sslCert", "", "Specify the path to the file containing the cert.pem file (filename must be included)")
	sslKeyPath  = flag.String("sslKey", "", "Specify the path to the file containing the key.pem file (filename must be included)")
)

type server struct {
	pb.UnsafeHumidityStorageServer
}

func main() {
	//handleHTTP()
	//handleGRPC()
	handleProtoBasedSocket()
}

func handleProtoBasedSocket() {
	localIp := "0.0.0.0"
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *port))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		log.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		log.Println("INFO: closing connection")
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	if n <= 0 {
		log.Println("no data received")
		return
	}

	var e pb.StoreHumidityRequest
	if err := proto.Unmarshal(buf[:n], &e); err != nil {
		log.Println("failed to unmarshal:", err)
		return
	}

	fmt.Printf("{DeviceID:%d, Humidity:%d}\n",
		e.GetDeviceId(),
		e.GetHumidity(),
	)

}

func handleGRPC() {
	localIp := "0.0.0.0"
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *port))
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

func handleHTTP() {
	http.HandleFunc("/receive", receiveHandler)
	log.Info("Starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Unable to start Server: ", err)
		return
	}
}

type Message struct {
	ID    int32 `json:"deviceid"`
	Value int32 `json:"value"`
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Received connection")
	if r.Method == "POST" {
		var message Message
		errJson := json.NewDecoder(r.Body).Decode(&message)
		if errJson != nil {
			http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error("Error while closing the http body: ", err.Error())
			}
		}(r.Body)

		forwardData(message.ID, message.Value)
	} else {
		http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
	}
}

func forwardData(id int32, value int32) {
	log.Info("Received ID:", id)
	log.Info("Received Value:", value)
	ForwardToHA(id, value)
	ForwardToPlantServer(id, value)
}

func (s server) StoreHumidityEntry(ctx context.Context, request *pb.StoreHumidityRequest) (*pb.StoreHumidityReply, error) {
	forwardData(request.DeviceId, request.Humidity)
	return &pb.StoreHumidityReply{}, nil
}
