package chat

import (
	"context"
	"daos_core/internal/data/models/chat"
	"daos_core/internal/infrastructure/db"
	"fmt"
)

type Repository interface {
	Save(ctx context.Context, input chat.CreateInput) error
}

type impl struct {
	postgres *db.Postgres
}

func NewRepository(postgres *db.Postgres) (Repository, error) {
	if postgres == nil {
		return nil, ErrPostgresArgument
	}
	return &impl{postgres: postgres}, nil
}

func (r *impl) Save(ctx context.Context, input chat.CreateInput) error {
	c, err := r.postgres.Pool.Exec(
		ctx,
		`INSERT INTO chat(amo_chat_id, user_id, instances_id) VALUES ($1, $2, $3);`,
		input.ConversationID,
		input.InstansceID,
		input.UserID,
	)
	if err != nil {
		return nil
	}

	if c.RowsAffected() == 0 {
		return fmt.Errorf("cant save chat")
	}

	return nil
}
