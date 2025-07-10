package userservice

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
)

type PostgresUserRepository struct {
	db *pgx.Conn
}

func NewPostgresUserRepository(db *pgx.Conn) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByLogin(login string) (*User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) Register(ctx context.Context, login string, email string, hashedPassword string) (int64, error) {
	var userID int64

	query := `
	INSERT INTO users (login, email, hashedPassword)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.db.QueryRow(ctx, query, login, email, hashedPassword).Scan(&userID)
	if err != nil {
		//TODO:
		//Я не знаю какие ошибки могут возникнуть при работе с БД, поэтому пока сделаю просто
		//Как я понял нужно проверять на то, что в Бд уже сущесвтуют пользователь или по крайней мере почта
		//Но я пока не знаю как это сделать
		//Как я понял это поля с UNIQUE, а также users_login_idx и users_email_idx. С помощью них можно будет проверить
		return -1, err
	}

	return userID, nil
}