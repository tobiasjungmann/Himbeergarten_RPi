module github.com/tobiasjungmann/Himbeergarten_RPi/testClient

go 1.21rc3

replace example.com/tobiasjungmann/Himbeergarten_RPi/server/proto => ../server/proto

require (
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.56.2
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/tobiasjungmann/Himbeergarten_RPi/server v0.0.0-20230726194455-e2982c509f05 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gorm.io/driver/sqlite v1.5.2 // indirect
	gorm.io/gorm v1.25.2 // indirect
)
