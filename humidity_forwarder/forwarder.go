package main

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	localAddress            = "0.0.0.0:12346"
	secretTokenPlantStorage = "secret_token"
	testTokenHomeAssistant  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI4MzczNDExNWFlNDc0ZGY4YjJiOGRlNWEzMDZkNTFkMCIsImlhdCI6MTY5MDY0NzM1MSwiZXhwIjoyMDA2MDA3MzUxfQ.RBXcYVaGhas-GPBt-04jE56TX1X50E7ypTJIKR-7zYQ"
)

func generateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretTokenPlantStorage))
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

		id := r.FormValue("deviceid")
		value := r.FormValue("value")
		log.Info("Received ID:", id)
		log.Info("Received Value:", value)
		forwardData(id, value, w)

	} else {
		http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
	}
}

func forwardData(id string, value string, w http.ResponseWriter) {
	idNumber, parsingIdError := parseToInt32(id, w)
	if parsingIdError != nil {
		return
	}
	valueNumber, ParsingValueError := parseToInt32(value, w)
	if ParsingValueError != nil {
		return
	}
	ForwardToHA(idNumber, valueNumber)
	ForwardToPlantServer(idNumber, valueNumber)
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
