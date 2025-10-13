package repository

import "context"

const (
	RemoveVerificationCodeQuery = `
	DELETE FROM email_verifications
	WHERE user_id = $1;
	`
)

func (r *PostgresUserRepository) RemoveCodeFromDB(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, RemoveVerificationCodeQuery, userID)

	return err
}