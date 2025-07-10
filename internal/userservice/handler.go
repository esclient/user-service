package userservice

import (
	"context"

	pb "github.com/esclient/user-service/api/userservice"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return nil, nil
}
