package telegram_request

// шлет фронт мне
type SendBotTokenDto struct {
	InstanceId int64  `json:"instanceId"`
	AccountId  int64  `json:"accountId"`
	BotToken   string `json:"botToken"`
}
