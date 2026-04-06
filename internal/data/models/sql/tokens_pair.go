package sqlmodels

import (
	"time"
)

// create table tokens_pair(
// 	id serial not null primary key,
// 	-- приходит от амо uuid
// 	user_id text not null,
// 	access_token text not null,
// 	refresh_token text not null,
// 	referer text not null,
// 	expires_at timestamp not null
// );

type TokensPair struct {
	Id           int       `db:"id"`
	UserId       string    `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	Referer      string    `db:"referer"`
	ExpiresAt    time.Time `db:"expires_at"`
	// InstanseLimit     int       `db:"instance_limit"`
	// InstanseAvailable int       `db:"instance_available"`
}
