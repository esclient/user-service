package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type DBCodeData struct {
	Code string
	CreatedAt time.Time
}

const (
	VerificationCodeLifetime = time.Minute * 5
)

const (
	GetCodeFromDBQuery = `
	SELECT code, created_at
	FROM email_verifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT 1;
	`
)

func (r *PostgresUserRepository) VerifyUser(ctx context.Context, userID int64, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, DBTimeout)
	defer cancel()

	dbCodeData, err := r.getCodeFromDB(ctx, userID)
    if err != nil {
        return false, err
    }

	codeValidation, err := validateCodeRepositoryLayer(code, dbCodeData)
	if err != nil {
		return codeValidation, err
	}

	if statusUpdErr := r.UpdateUserStatus(ctx, userID, UserActiveStatus); statusUpdErr != nil {
		log.Print("Ошибка в обновлении статуса пользователя")
		return false, statusUpdErr
	}

	if removeCodeErr := r.RemoveCodeFromDB(ctx, userID); removeCodeErr != nil {
		log.Print("Ошибка в удалении кода из email_verifications")
		return false, err
	}

	return true, nil
}

func (r *PostgresUserRepository) getCodeFromDB(ctx context.Context, userID int64) (DBCodeData, error) {
    var data DBCodeData
    
    err := r.db.QueryRow(ctx, GetCodeFromDBQuery, userID).Scan(&data.Code, &data.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return DBCodeData{}, ErrorRowDoesNotExist
        }
        return DBCodeData{}, ErrorQueryFailed
    }
    
    return data, nil
}

func validateCodeRepositoryLayer(code string, codeData DBCodeData) (bool, error) {
	if codeMismatchErr := isCodeMismatch(codeData.Code, code); codeMismatchErr != nil {
		return false, codeMismatchErr
	}

	if expiredCodeErr := isCodeExpired(codeData.CreatedAt); expiredCodeErr != nil {
		return false, expiredCodeErr
	}

	return true, nil
}

func isCodeExpired(createdAt time.Time) error {
	if time.Since(createdAt) > VerificationCodeLifetime {
		return ErrorCodeExpired
	}
	return nil
}

func isCodeMismatch(code string, codeFromDB string) error {
	if code != codeFromDB {
		return ErrorCodeMismatch
	}
	return nil
}