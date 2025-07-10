package userservice

import (
	"database/sql"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByLogin(login string) (*User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) Register(user *User) (int64, error) {
	return 0, nil
}
