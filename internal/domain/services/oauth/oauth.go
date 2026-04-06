package oauth

import (
	"context"
	oauth_models "daos_core/internal/data/models/oauth"
	account_repository "daos_core/internal/data/repositories/account"
	oauth_repository_impl "daos_core/internal/data/repositories/oauth"
	"daos_core/internal/domain/dto/oauth"
	"daos_core/internal/external/adapters/account"
	"daos_core/internal/external/adapters/channel"
	oauth_adapter "daos_core/internal/external/adapters/oauth"
	models "daos_core/internal/external/models/oauth"
	"fmt"
	"time"
)

type Service interface {
	UpdateAccessToken(ctx context.Context, dto oauth.RefreshTokensDTO) error
	SaveTokens(ctx context.Context, code string, referer string) error
}

type impl struct {
	OauthRepository   oauth_repository_impl.Repository
	AccountRepository account_repository.Repository
	OauthAdapter      oauth_adapter.Adapter
	AccountAdapter    account.Adapter
	ChannelAdapter    channel.Adapter
}

func NewOauthService(
	o_r oauth_repository_impl.Repository,
	a_r account_repository.Repository,
	o_a oauth_adapter.Adapter,
	a_a account.Adapter,
	c_a channel.Adapter,
) Service {
	return &impl{
		OauthRepository:   o_r,
		AccountRepository: a_r,
		OauthAdapter:      o_a,
		ChannelAdapter:    c_a,
		AccountAdapter:    a_a,
	}
}

func (s *impl) SaveTokens(ctx context.Context, code string, referer string) error {
	tokens, err := s.OauthAdapter.GetTokens(code, referer)
	if err != nil {
		return fmt.Errorf("OauthService: SaveTokens: %w", err)
	}

	if tokens.AccessToken == "" {
		return fmt.Errorf("OauthService: SaveTokens: access token cant be nil")
	}

	account, err := s.AccountAdapter.GetAccount(tokens.AccessToken, referer)
	if err != nil {
		return fmt.Errorf("OauthService: SaveTokens: %w", err)
	}

	if account == nil {
		return fmt.Errorf("OauthService: SaveTokens: account is nil")
	}

	// save account
	pk, err := s.AccountRepository.Create(ctx, *account)
	if err != nil {
		return fmt.Errorf("OauthService: SaveTokens: %w", err)
	}

	// amo tokens binding with accounts
	input := oauth_models.SaveInput{
		AccountPK:    pk,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Referer:      referer,
		ExpiredAt:    time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second),
	}

	err = s.OauthRepository.Create(ctx, input, referer)
	if err != nil {
		return fmt.Errorf("OauthService: SaveTokens: %w", err)
	}

	scopeId, err := s.ChannelAdapter.CreateChannel(account.AmojoID)
	if err != nil {
		return fmt.Errorf("OauthService: SaveTokens: %w", err)
	}

	if scopeId == nil {
		return fmt.Errorf("OauthService: SaveTokens: scope id is nil")
	}

	if err := s.AccountRepository.SaveScopeID(ctx, account.AmojoID, *scopeId); err != nil {
		return fmt.Errorf("OauthService: SaveTokens: %w", err)
	}
	return nil
}

func (s *impl) UpdateAccessToken(ctx context.Context, dto oauth.RefreshTokensDTO) error {
	if dto.Referer == "" {
		return fmt.Errorf("OauthService: UpdateAccessToken: referer was empty")
	}

	refresh, err := s.OauthRepository.GetRefreshToken(ctx, dto.Referer)
	if err != nil {
		return fmt.Errorf("OauthService: UpdateAccessToken: %w", err)
	}

	// fmt.Println("refresh - %s", token)

	data, err := s.OauthAdapter.RefreshTokens(models.RefreshInput{
		AccountID: dto.AccountID,
		Refresh:   refresh,
		Referer:   dto.Referer,
	})
	if err != nil {
		return fmt.Errorf("OauthService: UpdateAccessToken: %w", err)
	}

	// fmt.Printf("%+v\n", amoResponse)

	if data == nil {
		return fmt.Errorf("OauthService: UpdateAccessToken: data is nil")
	}

	expiresAt := time.Now().Add(time.Duration(data.ExpiresIn) * time.Second)
	input := oauth_models.UpdateInput{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiredAt:    expiresAt,
		Referer:      data.Referer,
	}

	err = s.OauthRepository.Update(ctx, input)
	if err != nil {
		return fmt.Errorf("OauthService: UpdateAccessToken: %w", err)
	}

	return nil
}
