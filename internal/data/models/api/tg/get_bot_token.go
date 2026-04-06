package tg_models

// что шлет виджет мне
type GetBotTokenRequest struct {
	UserId int64 `json:"user_id"`
	BotId  int64 `json:"bot_id"` // из нашей бд, а не тг
}

// шлю респонс прямиком из тг
