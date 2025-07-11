package handler

import (
	"context"
	"errors"
	"unicode"

	pb "github.com/esclient/user-service/api/userservice"

	"github.com/esclient/user-service/internal/userservice/service"
)

var (
	ErrorEmptyLogin = errors.New("login, field is empty")
	ErrorEmptyEmail = errors.New("email, field is empty")
	ErrorEmptyPassword = errors.New("password, field is empty")
	ErrorEmptyConfirmPassword = errors.New("confirm password, field is empty")
	ErrorConfirmPasswordMismatch = errors.New("the password field and confirm password do not match")

	ErrorCyrillicSymbolsLogin = errors.New("cyrillic symbols are not allowed in the login")
	ErrorCyrillicSymbolsEmail = errors.New("cyrillic symbols are not allowed in the email")
	ErrorCyrillicSymbolsPassword = errors.New("cyrillic symbols are not allowed in the password")
)

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
	err := validateRegisterRequest(req)
	if err != nil {
		return nil, err
	}

	login := req.Login
	email := req.Email
	password := req.Password

	userID, err := u.service.Register(ctx, login, email, password)
	if err != nil {
		return  nil, err
	}

	return &pb.RegisterUserResponse{UserId: userID}, nil
}

func validateRegisterRequest(req *pb.RegisterUserRequest) error {
	if hasCyrillic(req.Login) {
		return ErrorCyrillicSymbolsLogin
	}

	if hasCyrillic(req.Email) {
		return ErrorCyrillicSymbolsEmail
	}

	if hasCyrillic(req.Password) {
		return ErrorCyrillicSymbolsPassword
	}

	if isFieldEmpty(req.Login) {
		return ErrorEmptyLogin
	}
	if isFieldEmpty(req.Email) {
		return ErrorEmptyEmail
	}
	if isFieldEmpty(req.Password) {
		return ErrorEmptyPassword
	}
	if isFieldEmpty(req.ConfirmPassword) {
		return ErrorEmptyConfirmPassword
	}

	if !isPasswordConfirmMatch(req.Password, req.ConfirmPassword) {
		return ErrorConfirmPasswordMismatch
	}

	return nil
}

func isFieldEmpty(field string) bool {
	return len(field) == 0
}

func isPasswordConfirmMatch(password string, confirmPassword string) bool {
	return password == confirmPassword
}

func hasCyrillic(str string) bool {
	for _, r := range str {
        if unicode.Is(unicode.Cyrillic, r) {
            return true
        }
    }
    return false
}