package instance

type CreateDTO struct {
	AccountID string `json:"accountId" binding:"required"`
}
