package userservice

type User struct {
	ID          int64
	Login       string
	Email       string
	Password    string
}

type UserRepository interface {
	GetByLoginOrEmail(LoginOrEmail string) (*User, error)
	Register(user *User) (int64, error)
}
