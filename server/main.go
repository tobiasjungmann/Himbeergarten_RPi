package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasjungmann/Himbeergarten_RPi/server/models"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
)

var (
	port = flag.Int("port", 12346, "The server port")
)

type PlantStorage struct {
	pb.UnimplementedPlantStorageServer
	db *gorm.DB
}

const (
	secretToken = "secret_token"
)

func main() {
	dockerdb := flag.Bool("db", false, "Use MariaDB in Docker")
	flag.Parse()
	s := "test.db"
	if *dockerdb {
		s = "user:password@tcp(mariadb:3306)/mydatabase"
	}
	db, err := gorm.Open(sqlite.Open(s), &gorm.Config{})

	if err != nil {
		log.Fatalf("Terminating with error: %v", err)
		panic("failed to connect database")
	}
	errMigration := db.AutoMigrate(&models.Plant{}, &models.HumidityEntry{}, &models.ImageEntry{}, &models.Gpio{})
	if errMigration != nil {
		log.Fatalf("Unable to perform database migration. Terminating with error: %v", err)
		return
	}
	rpcServer(db)
}

func rpcServer(db *gorm.DB) {
	sslFlag := flag.Bool("ssl", false, "Enable SSL/TLS")
	flag.Parse()
	localIp := "0.0.0.0"
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIp, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *sslFlag {
		creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
		if err != nil {
			log.Fatalf("failed to load TLS certificates: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(tokenInterceptor)}
		log.Info("TLS is activated.")
	} else {
		opts = []grpc.ServerOption{grpc.UnaryInterceptor(tokenInterceptor)}
		log.Info("TLS is deactivated.")
	}

	s := grpc.NewServer(opts...)
	pb.RegisterPlantStorageServer(s, &PlantStorage{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func tokenInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
	}

	tokenString := authHeader[0]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretToken), nil
	})
	if err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return handler(ctx, req)
}
