package handler

import (
	"context"
	"errors"

	pb "github.com/esclient/user-service/api/userservice"

	"github.com/esclient/user-service/internal/userservice/service"
)

var ErrorValidFields = errors.New("login, email, password or confirm password field is empty")

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return nil, nil
}

func (u *UserHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	if IsValidRegistration(req) {
		return nil, ErrorValidFields
	}

	login := req.Login
	email := req.Email
	password := req.Password

	userID, err := u.service.Register(login, email, password)
	if err != nil {
		return  nil, err
	}

	return &pb.RegisterUserResponse{UserId: userID}, nil
}

func IsValidRegistration(req *pb.RegisterUserRequest) bool {
	return req.Login != "" &&
		req.Email != "" &&
		req.Password != "" &&
		req.ConfirmPassword != "" &&
		req.Password == req.ConfirmPassword
}