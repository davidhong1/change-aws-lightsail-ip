run:
	go run main.go -config config/config.yaml

dep:
	go mod download

build:
	go build -o bin/change-aws-lightsail-ip main.go
