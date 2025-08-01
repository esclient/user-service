package repository

import (
	"context"
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

	//TODO: Возможно стоит сделать проверку на существование пользователя
	
	return err
}