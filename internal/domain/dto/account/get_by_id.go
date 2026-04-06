package account

type TokensPairDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetByIDDTO struct {
	ID            int    `json:"id"`
	AccountID     int    `json:"account_id"`
	Subdomain     string `json:"subdomain"`
	Name          string `json:"name"`
	InstanceLimit int    `json:"instance_limit"`
	Language      string `json:"language"`
	Country       string `json:"country"`
}
