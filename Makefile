all: build run

build:
	go build -o ./bin/server.app ./cmd/main.go

run:
	DB_HOST=localhost DB_PORT=5432 DB_NAME=easy_normalization DB_USER=kolya59 DB_PASSWORD=12334566w REDIS_SERVER=localhost:6379  REDIS_DATABASE=0 REDIS_PASSWORD= ./bin/server.app


deploy-server:
	ssh -i "./.ssh/Macbook.pem" ubuntu@ec2-34-251-15-102.eu-west-1.compute.amazonaws.com

deploy-client:
	ssh -i "./.ssh/Macbook.pem" ubuntu@ec2-34-245-71-18.eu-west-1.compute.amazonaws.com

proto:
	protoc -I ./proto/ ./proto/Cars.proto --go_out=plugins=grpc:proto
