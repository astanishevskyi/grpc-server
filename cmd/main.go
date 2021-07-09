package main

import (
	"flag"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/configs"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/.env", "path to config file")
}

func main() {
	flag.Parse()
	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	storage := os.Getenv("STORAGE")
	config := configs.Config{
		BindAddr: port,
		Storage:  storage,
	}

	s := grpcserver.New(&config)
	if err := s.ConfigStorage(); err != nil {
		log.Fatal(err)
	}
	s.Run()
}
