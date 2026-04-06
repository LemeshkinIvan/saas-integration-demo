package chat_models

type CreateChatAmo struct {
	ConversationID string `json:"conversation_id"`
	User           User   `json:"user"`
}

type User struct {
	ID string `json:"id"`
}
