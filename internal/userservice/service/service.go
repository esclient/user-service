package service

import (
	"context"
	"os/user"

	crp "golang.org/x/crypto/bcrypt"

	"github.com/esclient/user-service/internal/userservice/repository"
)

const HashCost int = 14

type UserService struct {
	rep *repository.PostgresUserRepository
}

func NewUserService(rep *repository.PostgresUserRepository) *UserService {
	return &UserService{rep: rep}
}

func (s *UserService) Login(loginOrEmail string, password string) (*user.User, error) {
	return nil, nil
}

func (s *UserService) Register(ctx context.Context, login string, email string, password string) (int64, error) {
	hashedPassword, err := crp.GenerateFromPassword([]byte(password), HashCost)
	if err != nil {
		return -1, err
	}

	userID, err := s.rep.Register(ctx, login, email, string(hashedPassword))
	if err != nil {
		return -1, err
	}

	return userID, nil
}