dep:
	go mod download

run:
	go run main.go -config config/config.yaml

build:
	go build -o main main.go