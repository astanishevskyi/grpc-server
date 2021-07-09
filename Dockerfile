FROM golang:1.16-alpine

ENV config_path $config_path
RUN mkdir grpc-server
WORKDIR grpc-server
COPY . .

RUN go mod tidy
RUN echo $config_path
RUN go build ./cmd/main.go
EXPOSE 50051

