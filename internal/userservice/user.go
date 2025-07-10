package userservice

type User struct {
	ID          int64
	Login       string
	Email       string
	Password    string
}

type UserRepository interface {
	GetByLogin(login string) (*User, error)
	GetByEmail(email string) (*User, error)
	Register(user *User) (int64, error)
}
