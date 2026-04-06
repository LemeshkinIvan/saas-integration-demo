package max_models

type ChatResponse struct {
	Chats  []Chat `json:"chats"`
	Marker *int64 `json:"marker"`
}

type Chat struct {
	ChatId int64 `json:"chat_id"`
	// enum ChatType
	Type string `json:"type"`
	// enum ChatStatus
	Status            string         `json:"status"`
	Title             *string        `json:"title"`
	IconObj           *Icon          `json:"icon"`
	LastEventTime     int64          `json:"last_event_time"`
	ParticipantsCount int32          `json:"participants_count"`
	OwnerId           *int64         `json:"owner_id,omitempty"`
	ParticipantsObj   *Participants  `json:"participants,omitempty"`
	IsPublic          bool           `json:"is_public"`
	Link              *string        `json:"link,omitempty"`
	Description       *string        `json:"description"`
	DialogWithUser    *UserWithPhoto `json:"dialog_with_user,omitempty"`
	MessageCount      *int64         `json:"message_count,omitempty"`
	ChatMessageId     *string        `json:"chat_message_id,omitempty"`
	PinnedMessage     *Message       `json:"pinned_message,omitempty"`
}

type Participants struct{}
