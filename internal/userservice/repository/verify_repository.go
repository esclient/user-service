package repository

import (
	"context"
	"errors"
	"log"
	"time"
)

type DBCodeData struct {
	Code string
	CreatedAt string
}

var (
	ErrorRowsDoesNotExist = errors.New("db row does not exist")
	ErrorCodeExpired = errors.New("code from db expired")
	ErrorCodeMismatch = errors.New("code mismatch")
)

const (
	VerificationCodeLiftime = time.Minute * 5
)

const (
	GetCodeFromDBQuery = `
	SELECT code, created_at
	FROM email_verifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT 2;
	`
)

func (r *PostgresUserRepository) VerifyUser(ctx context.Context, userID int64, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, DBTimeout)
	defer cancel()

	rows, err := r.db.Query(ctx, GetCodeFromDBQuery, userID)
	if err != nil {
		return false, err
	}

	var DBCodeData DBCodeData
	if !rows.Next() {
		return  false, ErrorRowsDoesNotExist
	}
	rows.Scan(&DBCodeData.Code, &DBCodeData.CreatedAt)

	codeValidation, err := validateCodeRepositoryLayer(code, DBCodeData)
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

func validateCodeRepositoryLayer(code string, codeData DBCodeData) (bool, error) {
	if expiredCodeErr := isCodeExpired(codeData.CreatedAt); expiredCodeErr != nil {
		return false, expiredCodeErr
	}

	if codeMismatchErr := isCodeMismatch(codeData.Code, code); codeMismatchErr != nil {
		return false, codeMismatchErr
	}

	return true, nil
}

func isCodeExpired(createdAtAsString string) error {
	codeCreatedTime, formatTimeErr := time.Parse(TimeStampFormat, createdAtAsString)
	if formatTimeErr != nil {
		return formatTimeErr
	}

	if time.Since(codeCreatedTime) > VerificationCodeLiftime {
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