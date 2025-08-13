package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

const (
	UpdateUserStatusQuery = `
	UPDATE users
	SET status = $2
	WHERE id = $1;
	`
)

func (r *PostgresUserRepository) UpdateUserStatus(ctx context.Context, userID int64, status string) error {
	ctx, cancel := context.WithTimeout(ctx, DBTimeout)
	defer cancel()

	_, err := r.db.Exec(ctx, UpdateUserStatusQuery, userID, status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
            return ErrorUserNotFound
        }
	}
	
	return err
}