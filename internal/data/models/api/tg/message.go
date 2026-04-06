package tg_models

type Message struct {
	MessageId int64   `json:"message_id"`
	From      *User   `json:"from"`
	Date      int64   `json:"date"`
	Text      *string `json:"text"`
}

type User struct {
	Id        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	UserName  *string `json:"username"`
}
