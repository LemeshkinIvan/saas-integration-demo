package bot

import (
	"context"
	"daos_core/internal/data/models/bot"
	"daos_core/internal/infrastructure/db"
	"fmt"
)

type Repository interface {
	Delete(ctx context.Context, instanceID int64) error
	Create(ctx context.Context, token string) error
	CreateAndReturnPK(ctx context.Context, token string, botType int8) (int64, error)
	GetBotByInstanceID(ctx context.Context, instanceID int64) (*bot.BotModel, error)
	Update(ctx context.Context, newToken string, typeBot int8, tableID int64) error
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

func (r *impl) CreateAndReturnPK(ctx context.Context, token string, botType int8) (int64, error) {
	var pk int64

	err := r.postgres.Pool.QueryRow(
		ctx,
		`INSERT INTO bots_token (token, type_id) VALUES ($1, $2) Returning id`,
		token,
		botType,
	).Scan(&pk)

	if err != nil {
		return 0, fmt.Errorf("BotRepository: CreateAndReturnPK: %w", err)
	}

	return pk, nil
}

func (r *impl) GetBotByInstanceID(ctx context.Context, instanceID int64) (*bot.BotModel, error) {
	var bot bot.BotModel
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			Select *
			From bot_token
			WHERE instance_id = $1 
		`,
		instanceID,
	).Scan(
		&bot.ID,
		&bot.Token,
		&bot.InstanceID,
		&bot.CreatedAt,
		&bot.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("BotRepository: GetBotByInstanceID: %w", err)
	}

	return &bot, nil
}

func (r *impl) Delete(ctx context.Context, instanceID int64) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`
			DELETE FROM bots_token
			WHERE id IN (
				SELECT bot_token_id FROM instances WHERE id = $1 AND bot_token_id IS NOT NULL
			);
		`,
		instanceID,
	)
	if err != nil {
		return fmt.Errorf("BotRepository: Delete: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("BotRepository: Delete: rows affected 0")
	}

	return nil
}

func (r *impl) Update(ctx context.Context, newToken string, botType int8, tableID int64) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`UPDATE bots_token SET token = $1, type_id = $2 WHERE id = $3`,
		newToken,
		botType,
		tableID,
	)

	if err != nil {
		return fmt.Errorf("BotRepository: Update: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("BotRepository: Update: rows affected 0")
	}
	return nil
}

func (r *impl) Create(ctx context.Context, token string) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`INSERT INTO bots_token (token) VALUES ($1)`,
		token,
	)

	if err != nil {
		return fmt.Errorf("BotRepository: Save: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("BotRepository: Save: rows affected 0")
	}

	return nil
}
