package max_models

type Message struct {
	User      *User          `json:"sender,omitempty"`
	Recipient Recipient      `json:"recipient"`
	TimeStamp int64          `json:"timestamp"`
	Link      *LinkedMessage `json:"link,omitempty"`
	Body      *MessageBody   `json:"body"`
	Stat      *MessageStat   `json:"stat,omitempty"`
	Url       *string        `json:"url,omitempty"`
}

type MessageBody struct {
	Mid         string           `json:"mid"`
	Seq         int64            `json:"seq"`
	Text        *string          `json:"text"`
	Attachments *[]Attachment    `json:"attachments"`
	Markup      *[]MarkupElement `json:"markup,omitempty"`
}

type MessageStat struct {
	Views int `json:"views"`
}

type Attachment struct {
	Type    string                    `json:"type"`
	Payload *[]PhotoAttachmentPayload `json:"payload"`
}

type Recipient struct {
	ChatId int64 `json:"chat_id"`
	// enum ChatType
	ChatType string `json:"chat_type"`
	UserId   int64  `json:"user_id"`
}

type MarkupElement struct {
	Type   string `json:"type"`
	From   int32  `json:"from"`
	Length int32  `json:"length"`
}
