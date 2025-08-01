package repository

import "time"

const (
	UserPendingStatus = "PENDIG"
	UserActiveStatus = "ACTIVE"
)

const (
	DBTimeout = 5 * time.Second

	TimeStampFormat = "2006-01-02 15:04:05.999999"
)