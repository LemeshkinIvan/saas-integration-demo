package oauth

import (
	"context"
	db "daos_core/internal/infrastructure/db"
	"fmt"

	oauth_models "daos_core/internal/data/models/oauth"
)

type Repository interface {
	GetPK(ctx context.Context, referer string) (int, error)
	GetAccessToken(ctx context.Context, referer string) (string, error)
	GetRefreshToken(ctx context.Context, referer string) (string, error)
	Create(ctx context.Context, input oauth_models.SaveInput, referer string) error
	Update(ctx context.Context, input oauth_models.UpdateInput) error
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

func (r *impl) GetPK(ctx context.Context, referer string) (int, error) {
	var PK int

	err := r.postgres.Pool.QueryRow(ctx, `SELECT id FROM amo_auth WHERE referer = $1`, referer).Scan(&PK)
	if err != nil {
		return 0, fmt.Errorf("OauthRepository: GetPK: %w", err)
	}

	return PK, nil
}

func (r *impl) Create(ctx context.Context, input oauth_models.SaveInput, referer string) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`
			INSERT INTO amo_auth (referer, access_token, refresh_token, expires_at, account_id)
			VALUES ($1, $2, $3, $4, $5)
		`,
		referer,
		input.AccessToken,
		input.RefreshToken,
		input.ExpiredAt,
		input.AccountPK,
	)

	if err != nil {
		return fmt.Errorf("OauthRepository: Save:  %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("OauthRepository: Save: rows affected 0")
	}

	return nil
}

func (r *impl) GetAccessToken(ctx context.Context, referer string) (string, error) {
	var token string
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			SELECT access_token
			FROM amo_auth 
			WHERE referer = $1
		`,
		referer,
	).Scan(&token)

	if err != nil {
		return "", fmt.Errorf("OauthRepository: GetAccessToken: %w", err)
	}

	return token, nil
}

func (r *impl) GetRefreshToken(ctx context.Context, referer string) (string, error) {
	var result string
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			Select refresh_token
			FROM amo_auth
			WHERE referer = $1
		`,
		referer,
	).Scan(&result)

	if err != nil {
		return "", fmt.Errorf("OauthRepository: GetRefreshToken: %w", err)
	}

	return result, nil
}

func (r *impl) Update(ctx context.Context, input oauth_models.UpdateInput) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`
			UPDATE amo_auth
				SET access_token = $1,
					refresh_token = $2,
					expires_at = $3
			WHERE referer = $4
		`,
		input.AccessToken,
		input.RefreshToken,
		input.ExpiredAt,
		input.Referer,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("OauthRepository: Update: %w", err)
	}

	return nil
}
