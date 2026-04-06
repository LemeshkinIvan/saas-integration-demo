package user_dto

type GetUsers struct {
	Data EmbeddedUsers `json:"_embedded"`
}

type EmbeddedUsers struct {
	Users []User `json:"users"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
