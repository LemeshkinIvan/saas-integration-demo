package auth

import "time"

type SaveInput struct {
	AccountPK        int64
	AmoID            string
	Access           string
	RefreshHash      string
	AccessDuration   time.Duration
	RefreshExpiredAt time.Time
}

type UpdateInput struct {
	AccountID        string
	Access           string
	RefreshHash      string
	AccessDuration   time.Duration
	RefreshExpiredAt time.Time
}
