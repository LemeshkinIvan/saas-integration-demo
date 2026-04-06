package account

import (
	"context"
	accounts_response "daos_core/internal/data/models/accounts"
	models "daos_core/internal/data/models/oauth"
	"daos_core/internal/data/repositories/account"
	"daos_core/internal/data/repositories/oauth"
	dto "daos_core/internal/domain/dto/account"
	"fmt"
)

type Service interface {
	Get(ctx context.Context, accountID string) (*accounts_response.ShortAccountModel, error)
	Update(ctx context.Context, referer string, data models.AmoAccount) error
	GetTokens(ctx context.Context, accountID string) (*dto.TokensPairDTO, error)
	GetAmoAccessToken(ctx context.Context, referer string) (string, error)
	GetInstanceLimit(ctx context.Context, accountID string) (*dto.GetInstanceLimitDTO, error)
}

type impl struct {
	AccountRepository account.Repository
	OauthRepository   oauth.Repository
}

func NewService(a_r account.Repository, o_r oauth.Repository) Service {
	return &impl{
		AccountRepository: a_r,
		OauthRepository:   o_r,
	}
}

func (s *impl) GetInstanceLimit(ctx context.Context, accountID string) (*dto.GetInstanceLimitDTO, error) {
	if accountID == "" {
		return nil, fmt.Errorf("AccountInstance: GetInstanceLimit: accountId was empty")
	}

	data, err := s.AccountRepository.GetInstanceLimit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("AccountInstance: GetInstanceLimit: %w", err)
	}

	return &dto.GetInstanceLimitDTO{
		Limit: data,
	}, nil
}

func (s *impl) Get(ctx context.Context, accountID string) (*accounts_response.ShortAccountModel, error) {
	return nil, nil
}

func (s *impl) Update(ctx context.Context, referer string, data models.AmoAccount) error {
	return s.AccountRepository.Update(ctx, referer, data)
}

func (s *impl) GetTokens(ctx context.Context, accountID string) (*dto.TokensPairDTO, error) {
	return &dto.TokensPairDTO{}, nil
}

// TODO
func (s *impl) GetAmoAccessToken(ctx context.Context, referer string) (string, error) {
	token, err := s.OauthRepository.GetAccessToken(ctx, referer)
	if err != nil {
		return "", fmt.Errorf("AccountInstance: GetInstanceLimit: %w", err)
	}

	if token == "" {
		return "", fmt.Errorf("AccountInstance: GetInstanceLimit: %w", err)
	}

	return token, nil
}
