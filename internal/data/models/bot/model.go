package bot

import "time"

type BotModel struct {
	ID         int
	Token      string
	InstanceID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
