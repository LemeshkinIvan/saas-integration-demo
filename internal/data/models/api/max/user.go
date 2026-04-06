package max_models

type User struct {
	UserId           int64   `json:"user_id"`
	FirstName        string  `json:"first_name"`
	LastName         *string `json:"last_name"`
	Name             *string `json:"name,omitempty"`
	UserName         *string `json:"user_name"`
	IsBot            bool    `json:"is_bot"`
	LastActivityTime int64   `json:"last_activity_time"`
}

type UserWithPhoto struct {
	UserId           int64   `json:"user_id"`
	FirstName        string  `json:"first_name"`
	LastName         *string `json:"last_name"`
	Name             *string `json:"name,omitempty"`
	UserName         *string `json:"username"`
	IsBot            bool    `json:"is_bot"`
	LastActivityTime int64   `json:"last_activity_time"`
	Description      *string `json:"description,omitempty"`
	AvatarUrl        string  `json:"avatar_url,omitempty"`
	FullAvatarUrl    string  `json:"full_avatar_url,omitempty"`
}
