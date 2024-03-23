FROM golang:1.19

WORKDIR /app

COPY . .
RUN go mod download

EXPOSE 12348:12348

COPY *.go ./
RUN dir -s
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main
CMD ["./main","-portApi","12348","-ha","-proto","-ipHa","192.168.178.61"]