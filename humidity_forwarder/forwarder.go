package humidity_forwarder

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	pb "github.com/tobiasjungmann/Himbeergarten_RPi/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
	"time"
)

const (
	localAddress = "0.0.0.0:12346"
	secretToken  = "secret_token"
)

func generateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	http.HandleFunc("/receive", receiveHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Unable to start Server: ", err)
		return
	}
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data.", http.StatusInternalServerError)
			return
		}

		id, parsingIdError := parseToInt32(r.FormValue("deviceid"), w)
		if parsingIdError != nil {
			return
		}
		value, ParsingValueError := parseToInt32(r.FormValue("value"), w)
		if ParsingValueError != nil {
			return
		}

		// Print the received data to the console
		fmt.Println("Received ID:", id)
		fmt.Println("Received Value:", value)
		forwardToHA(id, value)
		forwardToPlantServer(id, value)
	} else {
		http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
	}
}

func parseToInt32(numberStr string, w http.ResponseWriter) (int32, error) {
	i, parsingError := strconv.ParseInt(numberStr, 10, 32)
	value := int32(i)
	if parsingError != nil {
		http.Error(w, "Invalid number format.", http.StatusBadRequest)
		return 0, parsingError
	}
	return value, nil
}

func forwardToPlantServer(id int32, value int32) {
	conn, err := grpc.Dial(localAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Info(err)
	}
	c := pb.NewPlantStorageClient(conn)
	s, _ := generateToken()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", s)

	_, errStore := c.StoreHumidityEntry(ctx, &pb.StoreHumidityRequest{RequestNumber: id, Humidity: value})

	if errStore != nil {
		log.Error(errStore.Error())
	}
}

func forwardToHA(id int32, number int32) {

}
