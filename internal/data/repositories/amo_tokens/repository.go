package amo_auth

import (
	"context"
	models "daos_core/internal/data/models/amo_tokens"
	"daos_core/internal/infrastructure/db"
	"fmt"
)

type Repository interface {
	GetAccessToken(ctx context.Context, accountID string) (string, error)
	GetReferer(ctx context.Context, accountID string) (string, error)

	GetAmoCredentials(ctx context.Context, accountID string) (*models.CredentialsModel, error)
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

func (r *impl) GetAmoCredentials(ctx context.Context, accountID string) (*models.CredentialsModel, error) {
	var referer string
	var access string

	err := r.postgres.Pool.QueryRow(
		ctx,
		`SELECT 
			a.access_token,
			a.referer
        FROM amo_auth a
        JOIN accounts ac ON ac.id = a.account_id
        WHERE ac.amo_id = $1`,
		accountID,
	).Scan(
		&access,
		&referer,
	)

	if err != nil {
		return nil, fmt.Errorf("AmoAuth: GetAmoCredentials: %w", err)
	}

	return &models.CredentialsModel{
		Referer:     referer,
		AccessToken: access,
	}, nil
}

func (r *impl) GetAccessToken(ctx context.Context, accountID string) (string, error) {
	const query = `
		SELECT t.access_token
        FROM accounts a
        JOIN amo_auth t ON t.id = a.amo_api_tokens_id
        WHERE a.amo_id = $1
	`
	var token string
	err := r.postgres.Pool.QueryRow(
		ctx,
		query,
		accountID,
	).Scan(&token)

	if err != nil {
		return "", fmt.Errorf("AmoTokensRepository: GetAccessToken: %w", err)
	}

	return token, nil
}

func (r *impl) GetReferer(ctx context.Context, accountID string) (string, error) {
	const query = `
		SELECT t.referer
        FROM accounts a
        JOIN amo_auth t ON t.id = a.amo_api_tokens_id
        WHERE a.account_id = $1
	`
	var referer string
	err := r.postgres.Pool.QueryRow(
		ctx,
		query,
		accountID,
	).Scan(&referer)

	if err != nil {
		return "", fmt.Errorf("AmoTokensRepository: GetReferer: %w", err)
	}

	return referer, nil
}
