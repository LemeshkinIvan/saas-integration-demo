package instance

type GetByIDDTO struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	BotToken   *string `json:"bot_token"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	Status     string  `json:"status"`
	SourceID   int     `json:"source_id"`
	PipelineID int     `json:"pipeline_id"`
	AmoErr     *string `json:"amo_err"`
}
