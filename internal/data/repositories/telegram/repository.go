package telegram

import (
	"context"
	db "daos_core/internal/infrastructure/db"
)

type Repository interface {
	GetToken(ctx context.Context, userID int64, botID int64) (string, error)
	GetStatusByAccount(ctx context.Context, accountID string) (string, error)
	GetStatusByInstance(ctx context.Context, instanceID int64) (string, error)
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

// на самом деле это instance
// id - account id, instance id
func (r *impl) GetStatusByInstance(ctx context.Context, instanceID int64) (string, error) {
	const getStatusQuery = `
		SELECT s.value
		FROM instances_status s
		JOIN instance i ON instance.status_id = s.id
		WHERE i.instanceId = $1
	`

	var status string
	err := r.postgres.Pool.QueryRow(ctx, getStatusQuery, instanceID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

// не трогай
func (r *impl) GetStatusByAccount(ctx context.Context, accountID string) (string, error) {
	const getStatusQuery = `
		SELECT s.value
		FROM instances_status s
		JOIN instance i ON instance.status_id = s.id
		WHERE i.instanceId = $1
	`

	var status string
	err := r.postgres.Pool.QueryRow(ctx, getStatusQuery, accountID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *impl) GetToken(ctx context.Context, userID int64, botID int64) (string, error) {
	const getBotQuery = `
		SELECT t.value
		FROM bot_user_tokens ut
		JOIN bot_token t ON ut.token_id == t.id
		WHERE id_user_tokens == $1 AND id_bot_token == $2
	`

	var token string
	err := r.postgres.Pool.QueryRow(
		ctx,
		getBotQuery,
		userID,
		botID,
	).Scan(
		&token,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
