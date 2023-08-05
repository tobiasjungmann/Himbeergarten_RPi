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
	port           = flag.Int("port", 12346, "The server port")
	dbUser         = flag.String("dbUser", "", "Specify the database user (e.g. user)")
	dbPwd          = flag.String("dbPwd", "", "Specify the database pwd (e.g. password)")
	dbType         = flag.String("dbType", "", "Specify the database type (e.g. mariadb)")
	dbPort         = flag.String("dbPort", "", "Specify the database (e.g. 3306)")
	dbDatabaseName = flag.String("dbdatabase", "", "Specify the database (e.g. mydatabase)")

	sslCertPath = flag.String("sslCert", "", "Specify the path to the file containing the cert.pem file (filename must be included)")
	sslKeyPath  = flag.String("sslKey", "", "Specify the path to the file containing the key.pem file (filename must be included)")
)

type PlantStorage struct {
	pb.UnimplementedPlantStorageServer
	db *gorm.DB
}

const (
	secretToken = "secret_token"
)

func main() {
	flag.Parse()
	s := "test.db"
	if len(*dbUser) > 0 {
		s = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPwd, *dbType, *dbPort, *dbDatabaseName) //"user:password@tcp(mariadb:3306)/mydatabase"
		log.Info("Using database connection String: %s", s)
	}
	db, err := gorm.Open(sqlite.Open(s), &gorm.Config{})

	if err != nil {
		log.Fatalf("Terminating with error: %v", err)
		panic("failed to connect database")
	}
	errMigration := db.AutoMigrate(&models.Plant{}, &models.HumidityEntry{}, &models.ImageEntry{}, &models.Sensor{})
	if errMigration != nil {
		log.Fatalf("Unable to perform database migration. Terminating with error: %v", err)
		return
	}
	rpcServer(db)
}

func rpcServer(db *gorm.DB) {
	go startSensorAPI(db)
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
		opts = []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(tokenInterceptor)}
		log.Info("TLS is activated.")
	} else {
		opts = []grpc.ServerOption{grpc.UnaryInterceptor(tokenInterceptor)}
		log.Info("TLS is deactivated.")
	}

	s := grpc.NewServer(opts...)
	pb.RegisterPlantStorageServer(s, &PlantStorage{db: db})
	log.Info("Plant Server listening at ", lis.Addr())
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
