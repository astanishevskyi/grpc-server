all: build run
.PHONY: build
config=
build:
	go build cmd/main.go
run:
	go run cmd/main.go
build-container:
	docker build -t grpc-server .
up:
	docker run grpc-server ./main $(config)
