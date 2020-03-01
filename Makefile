all: build run

build-all: build-mac build-win build-linux

build-mac: build-server-mac build-client-mac

build-server-mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/server.app ./cmd/server/main.go

build-client-mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/client.app ./cmd/client/main.go

build-win: build-server-win build-client-win

build-server-win:
	GOOS=windows GOARCH=amd64 go build -o ./bin/win/server.exe ./cmd/server/main.go

build-client-win:
	GOOS=windows GOARCH=amd64 go build -o ./bin/win/client.exe ./cmd/client/main.go

build-linux: build-server-linux build-client-linux

build-server-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux/server.bin ./cmd/server/main.go

build-client-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux/client.bin ./cmd/client/main.go

build-arm: build-server-arm build-client-arm

build-server-arm:
	GOOS=linux GOARCH=arm64 go build -o ./bin/arm/server.bin ./cmd/server/main.go

build-client-arm:
	GOOS=linux GOARCH=arm64 go build -o ./bin/arm/client.bin ./cmd/client/main.go


run: setup-env
	./bin/server.app

proto:
	protoc -I ./proto/ ./proto/Cars.proto --go_out=plugins=grpc:proto

setup-env:
	export `cat .env`