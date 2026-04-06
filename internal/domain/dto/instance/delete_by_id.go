package instance

type DeleteByIDDTO struct {
	SourceID  int64  `json:"source_id"`
	AccountID string `json:"account_id"`
}
