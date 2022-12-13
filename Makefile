run:
	go run main.go -config config/config.yaml

build:
	go build -o change-aws-lightsail-ip main.go
