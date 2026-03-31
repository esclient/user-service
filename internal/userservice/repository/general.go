package repository

import (
	"errors"
	"time"
)

const (
	UserPendingStatus = "PENDIG"
	UserActiveStatus = "ACTIVE"
)

const (
	DBTimeout = 5 * time.Second

	MaxPoolConns = 10
	MinPoolConns = 2
	MaxConnLifetime = time.Hour
	MaxConnIdleTime = 30 * time.Minute
)

var (
	ErrorQueryFailed     = errors.New("DB error code <0>") // db query failed

	ErrorLoginTaken 	 = errors.New("DB error code <3>") // login already taken
	ErrorEmailTaken 	 = errors.New("DB error code <4>") // email already taken

	ErrorUserNotFound    = errors.New("DB error code <3>") // user not found

	ErrorRowDoesNotExist = errors.New("DB error code <27>") // db row does not exist
	ErrorCodeExpired     = errors.New("DB error code <28>") // code from db expired
	ErrorCodeMismatch    = errors.New("DB error code <29>") // code mismatch
)