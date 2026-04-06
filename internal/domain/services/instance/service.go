package instance

import (
	"context"
	account "daos_core/internal/data/models/accounts"
	"daos_core/internal/data/models/amo_tokens"
	"daos_core/internal/data/models/instance"
	account_repository "daos_core/internal/data/repositories/account"
	amo_auth "daos_core/internal/data/repositories/amo_tokens"
	bot "daos_core/internal/data/repositories/bot"
	repository "daos_core/internal/data/repositories/instance"
	dto "daos_core/internal/domain/dto/instance"

	amo_source_client "daos_core/internal/external/adapters/source"
	"daos_core/internal/external/models/source"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	ListByAccountID(ctx context.Context, accountID string) (*dto.GetListDTO, error)
	GetByID(ctx context.Context, accountID string, instanceID int64) (*dto.GetByIDDTO, error)
	CountInstance(ctx context.Context, accountID string) (*dto.GetCountDTO, error)
	// возвращаю рыбу с undefined
	Create(ctx context.Context, accountID string) (*dto.GetInstanceDTO, error)
	Update(ctx context.Context, dto dto.UpdateDTO, accountID string) error
	Delete(ctx context.Context, accountID string, instanceID int64) error
}

type impl struct {
	InstanceRepo  repository.Repository
	BotTokensRepo bot.Repository
	AccountRepo   account_repository.Repository
	AmoTokensRepo amo_auth.Repository
	SourceAdapter amo_source_client.Adapter
}

func NewService(
	i_r repository.Repository,
	b_r bot.Repository,
	ac_r account_repository.Repository,
	a_r amo_auth.Repository,
	s_a amo_source_client.Adapter,
) Service {
	return &impl{
		InstanceRepo:  i_r,
		SourceAdapter: s_a,
		BotTokensRepo: b_r,
		AccountRepo:   ac_r,
		AmoTokensRepo: a_r,
	}
}

func (s *impl) GetByID(
	ctx context.Context,
	accountID string,
	instanceID int64,
) (*dto.GetByIDDTO, error) {
	if accountID == "" {
		return nil, fmt.Errorf("InstanceService: GetByID: accountId was empty")
	}

	item, err := s.InstanceRepo.GetByID(ctx, instanceID, accountID)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: GetByID: %w", err)
	}

	if item == nil {
		return nil, fmt.Errorf("InstanceService: GetByID: instance is nil")
	}

	return &dto.GetByIDDTO{
		ID:         item.ID,
		Name:       item.Name,
		CreatedAt:  item.CreatedAt,
		SourceID:   item.SourceID,
		PipelineID: item.PipelineID,
		Status:     item.Status,
		UpdatedAt:  item.UpdatedAt,
		BotToken:   item.BotToken,
	}, nil
}

func (s *impl) ListByAccountID(ctx context.Context, accountID string) (*dto.GetListDTO, error) {
	if accountID == "" {
		return nil, fmt.Errorf("InstanceService: ListByAccountID: accountID is empty")
	}

	list, err := s.InstanceRepo.ListByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: ListByAccountID: %w", err)
	}

	var data dto.GetListDTO
	if list == nil {
		data.Instances = []dto.GetInstanceDTO{}
	} else {
		for i := 0; i <= len(list)-1; i++ {
			data.Instances = append(data.Instances, dto.GetInstanceDTO{
				ID:        list[i].ID,
				Name:      list[i].Name,
				BotToken:  list[i].BotToken,
				CreatedAt: list[i].CreatedAt,
				Status:    list[i].Status,
			})
		}
	}

	limit, err := s.AccountRepo.GetInstanceLimit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: ListByAccountID: %w", err)
	}

	data.Limit = limit
	return &data, nil
}

func (s *impl) CountInstance(ctx context.Context, accountID string) (*dto.GetCountDTO, error) {
	if accountID == "" {
		return nil, fmt.Errorf("InstanceService: ListByAccountID: accountID is empty")
	}

	pk, err := s.AccountRepo.GetAccountPK(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: CountInstance: %w", err)
	}

	data, err := s.InstanceRepo.CountByAccountID(ctx, pk)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: CountInstance: %w", err)
	}

	return &dto.GetCountDTO{Count: data}, nil
}

// возвращаю рыбу с undefined
func (s *impl) Create(ctx context.Context, accountID string) (*dto.GetInstanceDTO, error) {
	if accountID == "" {
		return nil, fmt.Errorf("InstanceService: Create: accountId is empty")
	}

	auth, err := s.loadAmoAuthData(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: Create: %w", err)
	}

	meta, err := s.AccountRepo.GetAccountMeta(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: Create: %w", err)
	}

	if meta == nil {
		return nil, fmt.Errorf("InstanceService: Create: account meta is nil")
	}

	if err := s.checkInstanceLimit(ctx, *meta); err != nil {
		return nil, fmt.Errorf("InstanceService: Create: %w", err)
	}

	// создаю здесь uuid
	externalID := uuid.New().String()
	const defaultNameString = "Undefined"

	source, err := s.SourceAdapter.Create(
		source.AuthInput{
			Referer:     auth.Referer,
			AccessToken: auth.AccessToken,
		},
		source.CreateInput{
			Name:       defaultNameString,
			ExternalID: externalID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("InstanceService: Create: %w", err)
	}

	if err := s.validateSource(source); err != nil {
		return nil, fmt.Errorf("InstanceService: Create: %w", err)
	}

	// создаем полную сущность и отдаем short
	instance, err := s.InstanceRepo.Create(
		ctx,
		instance.CreateInput{
			AccountPK:  meta.PK,
			ExternalID: externalID,
			SourceID:   source.Embedded.Sources[0].ID,
			PipelineID: *source.Embedded.Sources[0].PipelineID,
		})
	if err != nil {
		return nil, fmt.Errorf("InstanceService: Create: %w", err)
	}

	if instance == nil {
		return nil, fmt.Errorf("InstanceService: Create: data is nil")
	}

	return &dto.GetInstanceDTO{
		ID:        instance.ID,
		Name:      instance.Name,
		BotToken:  instance.BotToken,
		CreatedAt: instance.CreatedAt,
		Status:    instance.Status,
	}, nil
}

func (s *impl) validateSource(source *source.GetSourcesResponseDTO) error {
	if source == nil {
		return fmt.Errorf("InstanceService: Create: source is nil")
	}

	if len(source.Embedded.Sources) == 0 {
		return fmt.Errorf("InstanceService: Create: source embedded is empty")
	}

	return nil
}

func (s *impl) checkInstanceLimit(ctx context.Context, meta account.MetaModel) error {
	currentNum, err := s.InstanceRepo.CountByAccountID(ctx, meta.PK)
	if err != nil {
		return fmt.Errorf("InstanceService: Create: %w", err)
	}

	if currentNum >= meta.InstanceLimit {
		return fmt.Errorf(
			"InstanceService: Create: Your limit(%d) has been exceeded for creating instances", meta.InstanceLimit,
		)
	}

	return nil
}

func (s *impl) loadAmoAuthData(ctx context.Context, accountID string) (*amo_tokens.CredentialsModel, error) {
	amoAuth, err := s.AmoTokensRepo.GetAmoCredentials(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if amoAuth.Referer == "" {
		return nil, fmt.Errorf("InstanceService: Create: referer not found")
	}

	if amoAuth.AccessToken == "" {
		return nil, fmt.Errorf("InstanceService: Create: referer not found")
	}

	return amoAuth, nil
}

func (s *impl) Update(
	ctx context.Context,
	data dto.UpdateDTO,
	accountID string,
) error {
	input := instance.UpdateInput{
		InstanceID: data.InstanceID,
		AccountID:  accountID,
		Name:       data.Name,
		PipelineID: *data.PipelineID,
		SourceID:   *data.SourceID,
	}

	botPK, err := s.InstanceRepo.Update(ctx, input)
	if err != nil {
		return fmt.Errorf("InstanceService: Update: %w", err)
	}

	if botPK != nil {
		err = s.BotTokensRepo.Update(ctx, data.Bot.Token, data.Bot.Type, *botPK)
		if err != nil {
			return fmt.Errorf("InstanceService: Update: %w", err)
		}
	} else {
		pk, err := s.BotTokensRepo.CreateAndReturnPK(ctx, data.Bot.Token, data.Bot.Type)
		if err != nil {
			return fmt.Errorf("InstanceService: Update: %w", err)
		}

		err = s.InstanceRepo.BindInstanceWithBotToken(ctx, data.InstanceID, pk)
		if err != nil {
			return fmt.Errorf("InstanceService: Update: %w", err)
		}
	}

	return nil
}

func (s *impl) Delete(
	ctx context.Context,
	accountID string,
	instanceID int64,
) error {
	sourceID, err := s.InstanceRepo.GetSourceIDByAccountID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("InstanceService: Delete: %w", err)
	}

	auth, err := s.loadAmoAuthData(ctx, accountID)
	if err != nil {
		return fmt.Errorf("InstanceService: Delete: %w", err)
	}

	fmt.Println(auth.Referer)
	fmt.Println(sourceID)

	if err := s.SourceAdapter.DeleteByID(
		source.AuthInput{
			Referer:     auth.Referer,
			AccessToken: auth.AccessToken,
		},
		sourceID,
	); err != nil {
		return fmt.Errorf("InstanceService: Delete: %w", err)
	}

	if err := s.BotTokensRepo.Delete(ctx, instanceID); err != nil {
		return fmt.Errorf("InstanceService: Delete: %w", err)
	}

	if err := s.InstanceRepo.DeleteByID(ctx, instanceID, accountID); err != nil {
		return fmt.Errorf("InstanceService: Delete: %w", err)
	}

	return nil
}

// func (s *impl) EnqueueInstanceServiceJob(job worker.Job) {
// 	worker.Pool.Enqueue(job)
// }
