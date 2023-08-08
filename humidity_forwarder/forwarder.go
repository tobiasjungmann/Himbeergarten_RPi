package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
)

var (
	apiPort           = flag.Int("portApi", 12348, "The forwarder port")
	portStorage       = flag.Int("portStorage", 12347, "The forwarder port")
	ipStorage         = flag.String("ipStorage", "0.0.0.0", "The forwarder port")
	ipHa              = flag.String("ipHa", "0.0.0.0", "The forwarder port")
	sslCertPath       = flag.String("sslCert", "", "Specify the path to the file containing the cert.pem file (filename must be included)")
	sslKeyPath        = flag.String("sslKey", "", "Specify the path to the file containing the key.pem file (filename must be included)")
	restReceiver      = flag.Bool("rest", false, "should the receiver accept rest requests to forward humidity data")
	bluetoothReceiver = flag.Bool("bluetooth", false, "should the receiver accept bluetooth connections to forward humidity data")
	protoReceiver     = flag.Bool("proto", false, "should the receiver accept requests to a socket in the proto format to forward humidity data")
	grpcReceiver      = flag.Bool("grpc", false, "should the receiver accept grpc requests to forward humidity data")
	haForwarder       = flag.Bool("ha", false, "should data also be forwarded to a local Home Assistant instance?")
	storageForwarder  = flag.Bool("storage", false, "should data also be forwarded to a plant storage instance?")
)

type server struct {
	pb.UnsafeHumidityStorageServer
}

func main() {
	flag.Parse()
	if !*restReceiver && !*bluetoothReceiver && !*protoReceiver && !*grpcReceiver {
		log.Fatalf("no API format selected. Set at least one.")
		return
	}
	if !*haForwarder && !*storageForwarder {
		log.Fatalf("No Forwarding option selected. Set at least one.")
		return
	}
	if *restReceiver {
		handleHTTP()
	}
	if *bluetoothReceiver {
		handleBluetooth()
	}
	if *protoReceiver {
		handleProtoBasedSocket()
	}
	if *grpcReceiver {
		handleGRPC()
	}
}

func forwardData(deviceId string, sensorId int32, humidity int32, humidityInPercent int32) {
	log.Printf("MAC: %s  Sensor: %d  Value: %d  ValueInPercent: %d\n", deviceId, sensorId, humidity, humidityInPercent)
	if *haForwarder {
		ForwardToHA(deviceId, sensorId, humidity, humidityInPercent)
	}
	if *storageForwarder {
		ForwardToPlantServer(deviceId, sensorId, humidity, humidityInPercent)
	}
}
