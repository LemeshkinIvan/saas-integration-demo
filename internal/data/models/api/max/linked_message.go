package max_models

type LinkedMessage struct {
	// MessageLinkType
	Type   string       `json:"type"`
	Sender *User        `json:"sender,omitempty"`
	ChatId int64        `json:"chat_id,omitempty"`
	Msg    *MessageBody `json:"message"`
}
