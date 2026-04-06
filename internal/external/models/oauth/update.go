package models

type UpdateTokensResponseDTO struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ServerTime   int    `json:"server_time"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateOutput struct {
	ExpiresIn    int
	AccessToken  string
	RefreshToken string
	Referer      string
}
