package grpcserver

import (
	"errors"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/configs"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/services"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/storage"
	pb "github.com/astanishevskyi/grpc-server/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	GRPCServer *grpc.Server
	config     *configs.Config
	Storage    models.UserService
}

func New(config *configs.Config) *Server {
	s := &Server{
		GRPCServer: grpc.NewServer(),
		config:     config,
	}
	if err := s.ConfigStorage(); err != nil {
		log.Fatal(err)
	}
	pb.RegisterUserServer(s.GRPCServer, &services.UserServer{DB: s.Storage})
	return s
}

func (s *Server) ConfigStorage() error {
	switch s.config.Storage {
	case "", "in-memory":
		log.Println("Storage is in-memory")
		inMemoryStorage := storage.NewInMemoryUserStorage()
		s.Storage = inMemoryStorage
		return nil
	case "redis":
		log.Println("Storage is redis")
		s.Storage = storage.NewRedisStorage()
		return nil
	case "mongo":
		log.Println("Storage is mongo")
		mongoStorage := storage.NewMongoStorage()
		s.Storage = mongoStorage
		return nil
	}
	return errors.New("no storage is set")
}

func (s *Server) Run() {
	log.Println("Server is running on " + s.config.BindAddr)
	lis, err := net.Listen("tcp", s.config.BindAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
