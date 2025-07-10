package userservice

type UserService struct {
	rep UserRepository
}

func NewUserService(rep UserRepository) *UserService {
	return &UserService{rep: rep}
}

func (s *UserService) Login(loginOrEmail string, password string) (*User, error) {
	return nil, nil
}
