package instance

import "time"

type Instance struct {
	ID         int64   `db:"id"`
	Name       string  `db:"name"`
	BotToken   *string `db:"bot_token"`
	CreatedAt  string  `db:"created_at"`
	UpdatedAt  *string `db:"updated_at"`
	Status     string  `db:"status"`
	SourceID   int     `db:"source_id"`
	PipelineID int     `db:"pipeline_id"`
}

type CreateInput struct {
	AccountPK  int64
	ExternalID string
	SourceID   int
	PipelineID int
}

type ShortInstance struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	BotToken  string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
	Status    string    `db:"status"`
}

type UpdateInput struct {
	InstanceID int64
	AccountID  string
	Name       string
	PipelineID int
	SourceID   int
}

// transfer between service and repository
type InputPair struct {
	Arg1 int
	Arg2 int
}
