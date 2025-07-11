package userservice

import (
	"context"
	"errors"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorLoginTaken = errors.New("login already taken")
	ErrorEmailTaken = errors.New("email already taken")
)

const RegisterUserQuery string = `
	INSERT INTO users (login, email, hashedPassword)
	VALUES ($1, $2, $3)
	RETURNING id
	`

type PostgresUserRepository struct {
	db *pgx.Conn
}

func NewDatabaseConnection(ctx context.Context, databaseURL string) (*pgx.Conn, error) {
	db, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
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

	err := r.db.QueryRow(ctx, RegisterUserQuery, login, email, hashedPassword).Scan(&userID)
	if err != nil {
		var pgxErr * pgconn.PgError
		if errors.As(err, &pgxErr) {
			switch pgxErr.ConstraintName {
				case "users_login_idx":
					return -1, ErrorLoginTaken
				case "users_email_idx":
					return -1, ErrorEmailTaken
			}
		}
		return  -1, err
	}

	return userID, nil
}