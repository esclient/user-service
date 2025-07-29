package repository

import (
	"context"
	"errors"
	"os/user"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorLoginTaken = errors.New("login already taken")
	ErrorEmailTaken = errors.New("email already taken")
)

const (
	RegisterUserQuery = `
	INSERT INTO users (login, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	CheckVerificationCodeExistsQuery = `SELECT EXISTS (SELECT 1 FROM email_verifications WHERE user_id = $1)`

	InsertVerificationCodeQuery = `
	INSERT INTO email_verifications (user_id, code, created_at)
	VALUES ($1, $2, $3)
	`

	UpdateVerificationCodeQuery = `
	UPDATE email_verifications SET code = $1 created_at = $2 WHERE user_id = $3
	`
)

const DBTimeout = 5 * time.Second

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

func (r *PostgresUserRepository) GetByLogin(login string) (*user.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*user.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) WriteVerificationCode(ctx context.Context, userID int64, verificationCode string) (error) {
	var exists bool

	ctx, cancel := context.WithTimeout(ctx, DBTimeout)
	defer cancel() 

	err := r.db.QueryRow(ctx, CheckVerificationCodeExistsQuery, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err := r.db.Exec(ctx, UpdateVerificationCodeQuery, verificationCode, time.Now(), userID)
		if err != nil {
			return  err
		}
	} else {
		_, err := r.db.Exec(ctx, InsertVerificationCodeQuery, userID, verificationCode, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresUserRepository) Register(ctx context.Context, login string, email string, hashedPassword string, verificationCode string) (int64, error) {
	var userID int64

	ctx, cancel := context.WithTimeout(ctx, DBTimeout)
	defer cancel()

	err := r.db.QueryRow(ctx, RegisterUserQuery, login, email, hashedPassword).Scan(&userID)
	if err != nil {
		var pgxErr *pgconn.PgError
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

	err = r.WriteVerificationCode(ctx, userID, verificationCode)
	if err != nil {
		return -1, err
	}

	return userID, nil
}