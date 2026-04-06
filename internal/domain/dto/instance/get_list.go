package instance

import "time"

type GetListDTO struct {
	Limit     int              `json:"limit" `
	Instances []GetInstanceDTO `json:"instances"`
}

type GetInstanceDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	BotToken  string    `json:"bot_token"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}
