package account

type GetByIDModel struct {
	ID            int    `db:"id"`
	AccountID     int    `db:"account_id"`
	Subdomain     string `db:"subdomain"`
	Name          string `db:"name"`
	InstanceLimit int    `db:"instance_limit"`
	AmojoID       string `db:"amojo_id"`
	ScopeID       string `db:"scope_id"`

	// Language string `db:"language"`

	// Country            string `db:"country"`
	// Currency           string `db:"currency"`
	// CurrencySymbol     string `db:"currency_symbol"`
	// IsHelpbotEnabled   bool   `db:"is_helpbot_enabled"`
	// IsTechnicalAccount bool   `db:"is_technical_account"`
}

type ShortAccountModel struct {
	ID            int    `db:"id"`
	AccountID     int    `db:"account_id"`
	Subdomain     string `db:"subdomain"`
	Name          string `db:"name"`
	InstanceLimit int    `db:"instance_limit"`
	Language      string `db:"language"`
	Country       string `db:"country"`
}
