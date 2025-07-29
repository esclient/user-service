package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"os/user"
	"time"
	"fmt"

	crp "golang.org/x/crypto/bcrypt"

	"github.com/esclient/user-service/internal/userservice/repository"
)

const (
	HashCost = 14
	VerificationCodeLiftime = time.Minute * 5
) 

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

	verificationCode, err := generateVerificationCode()
	if err != nil {
		return -1, err
	}

	userID, err := s.rep.Register(ctx, login, email, string(hashedPassword), verificationCode)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func generateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
    n, err := rand.Int(rand.Reader, max)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("%06d", n.Int64()), nil 
}