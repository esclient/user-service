package service

import "context"

func (s *UserService) VerifyUser(ctx context.Context, userID int64, code string) (bool, error) {
	verificationSuccess, err := s.rep.VerifyUser(ctx, userID, code)

	return  verificationSuccess, err
}