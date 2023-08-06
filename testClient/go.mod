module github.com/tobiasjungmann/Himbeergarten_RPi/testClient

go 1.19

replace example.com/tobiasjungmann/Himbeergarten_RPi/server/proto => ../server/proto

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/sirupsen/logrus v1.9.3
	github.com/tobiasjungmann/Himbeergarten_RPi/server v0.0.0-20230806151313-d8e2f2003e1f
	google.golang.org/grpc v1.56.2
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/image v0.1.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20230706204954-ccb25ca9f130 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230724170836-66ad5b6ff146 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230724170836-66ad5b6ff146 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gorm.io/driver/sqlite v1.5.2 // indirect
	gorm.io/gorm v1.25.2 // indirect
)
