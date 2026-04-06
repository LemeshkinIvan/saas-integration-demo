package auth

type GetTokensDTO struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
