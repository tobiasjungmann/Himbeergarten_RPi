FROM golang:1.19

WORKDIR /app

COPY . .
RUN go mod download

EXPOSE 12346:12346

COPY *.go ./
RUN dir -s
#RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main
CMD ["./main"]
