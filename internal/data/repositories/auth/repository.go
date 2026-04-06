package auth

import (
	"context"
	"daos_core/internal/data/models/auth"
	"daos_core/internal/infrastructure/db"
	"fmt"
	"time"
)

type Repository interface {
	GetRefresh(ctx context.Context, token string) (string, error)
	GetAccess(ctx context.Context, accountID string) (string, error)
	GetExpiredAt(ctx context.Context, accountID string) (*time.Time, error)
	Create(ctx context.Context, input auth.SaveInput) error
	Update(ctx context.Context, input auth.UpdateInput) error
	UpsertToken(ctx context.Context, input auth.SaveInput) error
}

type impl struct {
	Postgres *db.Postgres
	Cache    *db.RedisStorage
}

func NewRepository(p *db.Postgres, c *db.RedisStorage) (Repository, error) {
	if p == nil {
		return nil, ErrPostgresArgument
	}

	if c == nil {
		return nil, ErrCacheArgument
	}
	return &impl{
		Postgres: p,
		Cache:    c,
	}, nil
}

func (r *impl) UpsertToken(ctx context.Context, input auth.SaveInput) error {
	_, err := r.Postgres.Pool.Exec(ctx, `
		INSERT INTO tokens_pair(account_id, refresh_hash, expired_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (account_id) 
		DO UPDATE SET
			refresh_hash = EXCLUDED.refresh_hash,
			expired_at = EXCLUDED.expired_at
	`, input.AccountPK, input.RefreshHash, input.RefreshExpiredAt)

	if err != nil {
		return fmt.Errorf("AuthRepository: UpsertToken: %w", err)
	}

	// fmt.Println(input.AmoID)
	// fmt.Println(input.Access)

	if err := r.Cache.SetString(input.AmoID, input.Access, input.AccessDuration); err != nil {
		_, _ = r.Postgres.Pool.Exec(ctx,
			`DELETE FROM tokens_pair WHERE refresh_hash = $1`,
			input.RefreshHash,
		)
		return fmt.Errorf("AuthRepository: UpsertToken: %w", err)
	}

	return nil
}

func (r *impl) GetAccess(ctx context.Context, accountID string) (string, error) {
	token, err := r.Cache.GetString(accountID)
	if err != nil {
		return "", fmt.Errorf("AuthRepository: GetAccess: %w", err)
	}

	return token, nil
}

func (r *impl) GetRefresh(ctx context.Context, accountID string) (string, error) {
	var refrHash string
	err := r.Postgres.Pool.QueryRow(
		ctx,
		`
			Select t.token_hash
			From accounts a
			JOIN tokens_pair t ON t.id = a.tokens_pair_id 
			Where a.amo_id = $1
		`,
	).Scan(&refrHash)

	if err != nil {
		return "", fmt.Errorf("AuthRepository: GetRefresh: %w", err)
	}
	return refrHash, nil
}

func (r *impl) GetExpiredAt(ctx context.Context, accountID string) (*time.Time, error) {
	var exp *time.Time
	err := r.Postgres.Pool.QueryRow(
		ctx,
		`
			Select t.expired_at
			From accounts a
			JOIN tokens_pair t ON t.id = a.tokens_pair_id 
			Where a.amo_id = $1
		`,
		accountID,
	).Scan(&exp)

	if err != nil {
		return nil, fmt.Errorf("AuthRepository: Create: %w", err)
	}

	return exp, nil
}

func (r *impl) Create(ctx context.Context, input auth.SaveInput) error {
	res, err := r.Postgres.Pool.Exec(
		ctx,
		`
			WITH new_token AS (
				INSERT INTO tokens_pair(refresh_hash, expired_at)
				VALUES ($1, $2)
				RETURNING id
			)
			UPDATE accounts
			SET tokens_pair_id = new_token.id
			FROM new_token
			WHERE accounts.amo_id = $3;
		`,
		input.RefreshHash,
		input.RefreshExpiredAt,
		input.AccountPK,
	)

	if err != nil {
		return fmt.Errorf("AuthRepository: Create: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("AuthRepository: Create: rows affected 0")
	}

	if err := r.Cache.SetString(input.AmoID, input.Access, input.AccessDuration); err != nil {
		_, _ = r.Postgres.Pool.Exec(ctx,
			`DELETE FROM tokens_pair WHERE refresh_hash = $1`,
			input.RefreshHash,
		)
		return fmt.Errorf("AuthRepository: Create: %w", err)
	}

	return nil
}

func (r *impl) Update(ctx context.Context, input auth.UpdateInput) error {
	res, err := r.Postgres.Pool.Exec(
		ctx,
		`
			UPDATE tokens_pair tp
			SET refresh_hash = $1,
				expired_at   = $2
			FROM accounts a
			WHERE a.amo_id = $3
			AND a.token_pair_id = tp.id;
		`,
		input.RefreshHash,
		input.RefreshExpiredAt,
		input.AccountID,
	)

	if err != nil {
		return fmt.Errorf("AuthRepository: Update: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("AuthRepository: Update: rows affected 0")
	}

	if err := r.Cache.SetString(input.AccountID, input.Access, input.AccessDuration); err != nil {
		_, _ = r.Postgres.Pool.Exec(ctx,
			`DELETE FROM tokens_pair WHERE refresh_hash = $1`,
			input.RefreshHash,
		)
		return fmt.Errorf("AuthRepository: Update: %w", err)
	}

	return nil
}
