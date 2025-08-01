package handler

import (
	"context"
	"errors"
	"unicode"

	pb "github.com/esclient/user-service/api/userservice"
)

const CodeLength = 6

var (
	ErrorCodeEmpty = errors.New("confirmation code is empty")
	ErrorCodeLengthMismatch = errors.New("the code does not match the required length")
	ErrorCodeNotDigitable = errors.New("the code is not digitable")

	ErrorCodeUserIDNegative = errors.New("user ID is negative")
)

func (u *UserHandler) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	code := req.ConfirmationCode
	userID := req.UserId

	if codeErr := validateConfirmationCode(code); codeErr != nil {
		return nil, codeErr
	}

	if userIDErr := isUserIDNegative(userID); userIDErr != nil {
		return nil, userIDErr
	}

	verificationSuccess, err := u.service.VerifyUser(ctx, userID, code)

	return &pb.VerifyUserResponse{IsVerified: verificationSuccess}, err
}

func validateConfirmationCode(code string) error {
	var err error = nil

	if err = isCodeEmpty(code); err != nil {
		return err
	}
	if err = isCodeLengthMismatch(code); err != nil {
		return err
	}
	if err = isCodeNotDigitable(code); err != nil {
		return err
	}

	return nil	
}

func isCodeEmpty(code string) error {
	if (code == "") {
		return ErrorCodeEmpty
	}
	return nil
}

func isCodeLengthMismatch(code string) error {
	if (len(code) != CodeLength) {
		return ErrorCodeLengthMismatch
	}
	return nil
}

func isCodeNotDigitable(code string) error {
	for _, r := range code {
		if !unicode.IsDigit(r) {
			return ErrorCodeNotDigitable
		}
	}

	return nil
}

func isUserIDNegative(userID int64) error {
	if userID < 0 {
		return ErrorCodeUserIDNegative
	}
	return nil
}