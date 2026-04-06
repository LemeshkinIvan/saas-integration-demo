package sqlmodels

type Accounts struct {
	Id             int64  `db:"id"`
	TokensPairFK   *int64 `db:"tokens_pair_id"`
	AmoApiTokensFK int64  `db:"amo_api_tokens_id"`
	// amo
	AccountId int64  `db:"account_id"`
	Subdomain string `db:"subdomain"`
	Name      string `db:"name"`

	InstanceLimit int64 `db:"instance_limit"`
}
