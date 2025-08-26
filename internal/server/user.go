package server

import (
	"GRPCProject/internal/core"
	pb "GRPCProject/proto"
	"context"
)

type UserRepository interface {
	GetById(ctx context.Context, id string) (*core.User, error)
}

type UserServer struct {
	pb.UserServiceServer
	userRepository UserRepository
}

func NewUserServer(repository UserRepository) *UserServer {
	return &UserServer{
		userRepository: repository,
	}
}

func (server *UserServer) GetUser(ctx context.Context, request *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := server.userRepository.GetById(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &core.User{
			ID:        0,
			FirstName: "",
			LastName:  "",
		}
	}

	return &pb.UserResponse{Name: user.FirstName}, nil
}
