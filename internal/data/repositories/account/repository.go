package account

import (
	"context"
	account "daos_core/internal/data/models/accounts"
	m "daos_core/internal/data/models/oauth"
	"daos_core/internal/infrastructure/db"
	"fmt"
)

type Repository interface {
	Get(ctx context.Context, accountID string) (*account.GetByIDModel, error)
	Create(ctx context.Context, data m.AmoAccountOutput) (int64, error) // account
	Update(ctx context.Context, referer string, input m.AmoAccount) error

	GetInstanceLimit(ctx context.Context, accountID string) (int, error)
	GetAccountPK(ctx context.Context, accountID string) (int64, error)
	GetScopeID(ctx context.Context, accountID string) (string, error)
	GetScopeIDByInstance(ctx context.Context, instanceID int) (string, error)
	GetAccountMeta(ctx context.Context, accountID string) (*account.MetaModel, error)
	SaveScopeID(ctx context.Context, amojoID string, scopeID string) error
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

func (r *impl) SaveScopeID(ctx context.Context, amojoID string, scopeID string) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`
			UPDATE accounts
			SET scope_id = $2
			WHERE amojo_id = $1
		`,
		amojoID,
		scopeID,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("OauthRepository: SaveScopeID: rows affected 0")
	}

	return nil
}

func (r *impl) GetAccountMeta(ctx context.Context, accountID string) (*account.MetaModel, error) {
	var PK int64
	var instanceLimit int

	err := r.postgres.Pool.QueryRow(
		ctx,
		`Select id, instance_limit From accounts Where amo_id = $1`,
		accountID,
	).Scan(
		&PK,
		&instanceLimit,
	)

	if err != nil {
		return nil, fmt.Errorf("AccountRepository: GetAccountMeta: %w", err)
	}

	return &account.MetaModel{
		PK:            PK,
		InstanceLimit: instanceLimit,
	}, nil
}

func (r *impl) Create(ctx context.Context, data m.AmoAccountOutput) (int64, error) {
	var pk int64
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			INSERT INTO accounts( 
				amo_id, 
				name, 
				subdomain, 
				amojo_id
			) 
			VALUES($1, $2, $3, $4)
			RETURNING id;
		`,
		data.AmoID,
		data.Name,
		data.Subdomain,
		data.AmojoID,
	).Scan(&pk)

	if err != nil {
		return 0, fmt.Errorf("AccountRepository: Create: %w", err)
	}

	return pk, nil
}

func (r *impl) Get(ctx context.Context, accountID string) (*account.GetByIDModel, error) {
	var result account.GetByIDModel
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			SELECT 
				id, 
				amo_id, 
				subdomain, 
				name, 
				instance_limit
			FROM accounts
			WHERE amo_id = $1
		`,
		accountID,
	).Scan(
		&result.ID,
		&result.AccountID,
		&result.Subdomain,
		&result.Name,
		&result.InstanceLimit,
		// &result.Language,
		// &result.Country,
		// &result.Currency,
		// &result.CurrencySymbol,
		// &result.IsHelpbotEnabled,
		// &result.IsTechnicalAccount,
	)

	if err != nil {
		return nil, fmt.Errorf("AccountRepository: Get: %w", err)
	}

	return &result, nil
}

func (r *impl) Update(
	ctx context.Context,
	referer string,
	input m.AmoAccount,
) error {
	tx, err := r.postgres.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepository: Update: %w", err)
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(
		ctx,
		`
			UPDATE accounts a
			SET
				subdomain = $1,
				name = $2,
				amojo_id = $3
			FROM amo_api_tokens t
			WHERE a.token_id = t.id   -- внешний ключ
			AND t.referer = $4; 
		`,
		input.Subdomain,
		input.Name,
		input.AmojoID,
		referer,
	)

	if err != nil {
		return fmt.Errorf("AccountRepository: Update: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("AccountRepository: Update: row affected 0")
	}

	return tx.Commit(ctx)
}

func (r *impl) GetAccountPK(
	ctx context.Context,
	accountID string,
) (int64, error) {
	var key int64
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			SELECT id FROM accounts Where amo_id = $1
		`,
		accountID,
	).Scan(&key)

	if err != nil {
		return 0, fmt.Errorf("AccountRepository: GetAccountPK: %w", err)
	}

	return key, nil
}

func (r *impl) GetInstanceLimit(
	ctx context.Context,
	accountID string,
) (int, error) {
	var count int
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			Select instance_limit
			From accounts
			Where amo_id = $1
		`,
		accountID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("AccountRepository: GetInstanceLimit: %w", err)
	}

	fmt.Printf("limit:%d\n", count)
	return count, nil
}

func (r *impl) GetScopeID(ctx context.Context, accountID string) (string, error) {
	var scopeId string
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			SELECT scope_id
			FROM accounts
			WHERE amo_id = $1
		`,
		accountID,
	).Scan(&scopeId)

	if err != nil {
		return "", fmt.Errorf("AccountRepository: GetScopeID: %w", err)
	}

	return scopeId, nil
}

func (r *impl) GetScopeIDByInstance(ctx context.Context, instanceID int) (string, error) {
	var scopeID string
	var accountID int

	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			SELECT
				i.account_id,
				a.scope_id
			FROM instance i
			JOIN accounts a ON a.amo_id = i.account_id
			WHERE i.id = $1
		`,
		instanceID,
	).Scan(&accountID, &scopeID)

	if err != nil {
		return "", fmt.Errorf("AccountRepository: GetScopeIDByInstance: %w", err)
	}

	return scopeID, nil
}
