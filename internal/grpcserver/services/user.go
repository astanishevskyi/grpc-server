package services

import (
	"context"
	"github.com/astanishevskyi/grpc-server/internal/grpcserver/models"
	pb "github.com/astanishevskyi/grpc-server/pkg/api"
	"log"
)

type UserServer struct {
	pb.UnimplementedUserServer
	DB models.UserService
}

func (s *UserServer) GetUser(_ context.Context, in *pb.UserId) (*pb.UserObject, error) {
	log.Printf("Received: %v", in.GetId())
	resp, err := s.DB.Retrieve(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UserObject{Id: resp.ID, Age: uint32(resp.Age), Name: resp.Name, Email: resp.Email}, nil
}

func (s *UserServer) GetUsers(_ *pb.NoneObject, stream pb.User_GetUsersServer) error {
	for _, user := range s.DB.GetAll() {
		if err := stream.Send(&pb.UserObject{Id: user.ID, Age: uint32(user.Age), Name: user.Name, Email: user.Email}); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserServer) CreateUser(_ context.Context, in *pb.NewUser) (*pb.UserObject, error) {
	user, err := s.DB.Add(in.GetName(), in.GetEmail(), uint8(in.GetAge()))
	if err != nil {
		return nil, err
	}
	return &pb.UserObject{Id: user.ID, Age: uint32(user.Age), Name: user.Name, Email: user.Email}, nil
}

func (s *UserServer) UpdateUser(_ context.Context, in *pb.UserObject) (*pb.UserObject, error) {
	log.Printf("Received: %v", in.GetId())
	user, err := s.DB.Update(in.GetId(), in.GetName(), in.GetEmail(), uint8(in.GetAge()))
	if err != nil {
		return nil, err
	}
	return &pb.UserObject{Id: user.ID, Name: user.Name, Email: user.Email, Age: uint32(user.Age)}, nil
}

func (s *UserServer) DeleteUser(_ context.Context, in *pb.UserId) (*pb.UserId, error) {
	log.Printf("Received: %v", in.GetId())
	userID, err := s.DB.Remove(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UserId{Id: userID}, nil
}
