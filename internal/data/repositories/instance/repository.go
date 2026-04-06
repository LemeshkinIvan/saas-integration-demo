package instance

import (
	"context"
	instance "daos_core/internal/data/models/instance"
	source "daos_core/internal/external/models/source"
	"daos_core/internal/infrastructure/db"
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	ListByAccountID(ctx context.Context, accountID string) ([]instance.ShortInstance, error)
	GetByID(ctx context.Context, instanceID int64, accountID string) (*instance.Instance, error)
	GetSourceIDByAccountID(ctx context.Context, accountID string) (int, error)

	Create(ctx context.Context, input instance.CreateInput) (*instance.ShortInstance, error)
	CreateShortByAccountPK(ctx context.Context, accountPK int64) (*instance.ShortInstance, error)

	Update(ctx context.Context, input instance.UpdateInput) (*int64, error)
	UpdateStatus(ctx context.Context, instanceID int64, statusID int64) error
	UpdateBySource(ctx context.Context, instanceID int64, source source.SourceResponse, updatedAt time.Time) error

	DeleteByID(ctx context.Context, instanceID int64, accountID string) error

	CountByAccountID(ctx context.Context, accountPK int64) (int, error)
	BindInstanceWithBotToken(ctx context.Context, instanceID int64, bot int64) error
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

func (r *impl) GetSourceIDByAccountID(ctx context.Context, accountID string) (int, error) {
	var ID int
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			Select source_id 
			FROM instances i 
			JOIN accounts a ON a.id = i.account_id 
			WHERE a.amo_id = $1
		`,
		accountID,
	).Scan(&ID)

	if err != nil {
		return 0, fmt.Errorf("InstanceRepository: Create: %w", err)
	}
	return ID, nil
}

func (r *impl) Create(ctx context.Context, input instance.CreateInput) (*instance.ShortInstance, error) {
	var item instance.ShortInstance

	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			INSERT INTO instances (account_id, source_id, external_id, pipeline_id)
			VALUES ($1)
			RETURNING 
				instances.id,
				instances.name,
				instances.created_at,
				(SELECT value 
				FROM instances_status
				WHERE id = instances.status_id
			) AS status_value;
		`,
	).Scan(
		&item.ID,
		&item.Name,
		&item.CreatedAt,
		&item.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("InstanceRepository: Create: %w", err)
	}

	return &item, nil
}

func (r *impl) BindInstanceWithBotToken(ctx context.Context, instanceID int64, bot int64) error {
	res, err := r.postgres.Pool.Exec(
		ctx,
		`UPDATE instances SET bot_token_id = $1 Where id = $2`,
		bot,
		instanceID,
	)

	if err != nil {
		return fmt.Errorf("InstanceRepository: BindInstanceWithBotToken: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("InstanceRepository: BindInstanceWithBotToken: rows affected 0")
	}

	return nil
}

// +
func (r *impl) GetByID(ctx context.Context, instanceID int64, accountID string) (*instance.Instance, error) {
	if instanceID <= 0 {
		return nil, fmt.Errorf("InstanceRepository: GetByID: invalid instanceId")
	}

	var instance instance.Instance
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			SELECT 
				i.id, 
				b.token as bot_token,
				i.name,
				TO_CHAR(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at,
				TO_CHAR(i.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at,
				s.value as status,
				i.source_id,
				i.pipeline_id,
			FROM instances i
			JOIN accounts a ON a.id = i.account_id
			JOIN bots_token b ON b.id = i.bot_token_id
			JOIN instances_status s ON s.id = i.status_id
			WHERE a.amo_id = $1 and i.id = $2
		`,
		accountID,
		instanceID,
	).Scan(
		&instance.ID,
		&instance.BotToken,
		&instance.Name,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.Status,
		&instance.SourceID,
		&instance.PipelineID,
	)

	if err != nil {
		return nil, fmt.Errorf("InstanceRepository: GetByID: %w", err)
	}

	return &instance, nil
}

// +
func (r *impl) UpdateStatus(ctx context.Context, instanceID int64, statusID int64) error {
	if instanceID <= 0 {
		return fmt.Errorf("InstanceRepository: UpdateStatus: invalid instanceId")
	}

	if statusID <= 0 {
		return fmt.Errorf("InstanceRepository: DeleteByID: invalid statusID")
	}

	tx, err := r.postgres.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("InstanceRepository: UpdateStatus: %w", err)
	}

	defer tx.Rollback(ctx)

	res, err := tx.Exec(
		ctx,
		`
			UPDATE instances
			SET status_id = $1
			WHERE id = $2
		`,
		statusID,
		instanceID,
	)

	if err != nil {
		return fmt.Errorf("InstanceRepository: UpdateStatus: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("InstanceRepository: UpdateStatus: can't update")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("InstanceRepository: UpdateStatus: %w", err)
	}

	return nil
}

// +
func (r *impl) ListByAccountID(ctx context.Context, accountID string) ([]instance.ShortInstance, error) {
	var instances []instance.ShortInstance
	rows, err := r.postgres.Pool.Query(
		ctx,
		`
			SELECT 
					i.id,
					b.token AS token,
					i.name,
					i.created_at,
					s.value AS status
			FROM accounts a
			JOIN instances i ON i.account_id = a.id
			JOIN instances_status s ON i.status_id = s.id
			LEFT JOIN bots_token b ON i.bot_token_id = b.id
			WHERE a.amo_id = $1
		`,
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("InstanceRepository: ListByAccountID: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var item instance.ShortInstance
		var token sql.NullString

		err := rows.Scan(
			&item.ID,
			&token,
			&item.Name,
			&item.CreatedAt,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("InstanceRepository: ListByAccountID: %w", err)
		}

		if token.Valid {
			item.BotToken = token.String
		}

		instances = append(instances, item)
	}

	fmt.Printf("len: %d", len(instances))
	return instances, nil
}

// +
func (r *impl) CreateShortByAccountPK(ctx context.Context, accountPK int64) (*instance.ShortInstance, error) {
	if accountPK <= 0 {
		return nil, fmt.Errorf("InstanceRepository: CreateShortByAccountPK: invalid account PK")
	}

	var instance instance.ShortInstance
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			INSERT INTO instances (account_id)
			VALUES ($1)
			RETURNING 
				instances.id,
				instances.name,
				instances.created_at,
				(SELECT value 
				FROM instances_status
				WHERE id = instances.status_id
			) AS status_value;
		`,
		accountPK,
	).Scan(
		&instance.ID,
		&instance.Name,
		&instance.CreatedAt,
		&instance.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("InstanceRepository: CreateByAccountPK: %w", err)
	}

	return &instance, nil
}

func (r *impl) UpdateBySource(
	ctx context.Context,
	instanceID int64,
	source source.SourceResponse,
	updatedAt time.Time,
) error {
	if instanceID <= 0 {
		return fmt.Errorf("InstanceRepository: UpdateBySource: invalid instanceId")
	}

	res, err := r.postgres.Pool.Exec(
		ctx,
		`
			UPDATE instances
			SET external_id = $1, source_id = $2, pipeline_id = $3, updated_at = $4
			WHERE id = $5
		`,
		source.ExternalID,
		source.ID,
		source.PipelineID,
		updatedAt,
		instanceID,
	)

	if err != nil {
		return fmt.Errorf("InstanceRepository: UpdateBySource: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("InstanceRepository: UpdateBySource: rows affecetd 0")
	}

	return nil
}

// +
func (r *impl) Update(ctx context.Context, input instance.UpdateInput) (*int64, error) {
	if input.InstanceID <= 0 {
		return nil, fmt.Errorf("InstanceRepository: Update: invalid instanceId")
	}

	var botTokenID sql.NullInt64
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
			UPDATE instances
			SET name = $2,
				updated_at = NOW(),
				pipeline_id = $4
			WHERE id = $1
			AND account_id = (SELECT id FROM accounts WHERE amo_id = $3)
			RETURNING bot_token_id;
		`,
		input.InstanceID,
		input.Name,
		input.AccountID,
		input.PipelineID,
	).Scan(&botTokenID)

	if err != nil {
		return nil, fmt.Errorf("InstanceRepository: Update: %w", err)
	}

	if botTokenID.Valid {
		return &botTokenID.Int64, nil
	}
	return nil, nil
}

// +
// accountID - id from platform
func (r *impl) CountByAccountID(
	ctx context.Context,
	accountPK int64,
) (int, error) {
	if accountPK <= 0 {
		return 0, fmt.Errorf("InstanceRepository: CountByAccountID: invalid accountPK")
	}

	var num int
	err := r.postgres.Pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM instances i
		JOIN accounts a ON i.account_id = a.id
		WHERE a.id = $1
		`,
		accountPK,
	).Scan(&num)

	if err != nil {
		return 0, fmt.Errorf("InstanceRepository: CountByAccountID: %w", err)
	}

	return num, nil
}

// +
func (r *impl) DeleteByID(
	ctx context.Context,
	instanceID int64,
	accountID string,
) error {
	if instanceID <= 0 {
		return fmt.Errorf("InstanceRepository: DeleteByID: invalid instanceId")
	}

	res, err := r.postgres.Pool.Exec(
		ctx,
		`
			DELETE FROM instances i
			USING accounts a
			WHERE i.account_id = a.id   
				AND a.amo_id = $2   
				AND i.id = $1
			RETURNING i.id;  
		`,
		instanceID,
		accountID,
	)

	if res.RowsAffected() == 0 {
		return fmt.Errorf("InstanceRepository: DeleteByID: rows affected 0")
	}

	if err != nil {
		return fmt.Errorf("InstanceRepository: DeleteByID: %w", err)
	}

	return nil
}
