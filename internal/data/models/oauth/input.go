package oauth_models

import "time"

type UpdateInput struct {
	AccessToken  string
	RefreshToken string
	Referer      string
	ExpiredAt    time.Time
}

type SaveInput struct {
	AccountPK    int64
	AccessToken  string
	RefreshToken string
	Referer      string
	ExpiredAt    time.Time
}
