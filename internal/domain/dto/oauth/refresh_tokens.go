package oauth

type RefreshTokensDTO struct {
	AccountID int    `json:"account_id"`
	Referer   string `json:"referer"`
}
