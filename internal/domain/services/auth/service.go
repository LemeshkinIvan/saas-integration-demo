package auth

import (
	"context"
	"crypto/sha256"
	"daos_core/internal/constants"
	models "daos_core/internal/data/models/auth"
	"daos_core/internal/data/repositories/account"
	"daos_core/internal/data/repositories/auth"
	dto "daos_core/internal/domain/dto/auth"
	"daos_core/internal/utils/jwt"
	"encoding/hex"
	"fmt"
	"time"
)

type Service interface {
	Login(ctx context.Context, accountID string) (*dto.GetTokensDTO, error)
	ValidateAccess(ctx context.Context, token string) (string, error)
	RefreshToken(ctx context.Context, token string) (*dto.GetTokensDTO, error)
}

type impl struct {
	AuthConfig  *constants.AuthConfig
	AccountRepo account.Repository
	AuthRepo    auth.Repository
	JWT         jwt.JWTUtil
}

func NewService(cfg *constants.AuthConfig, au auth.Repository, acc account.Repository, jwt jwt.JWTUtil) (Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("AuthService: Login: cfg is nil")
	}

	if jwt == nil {
		return nil, fmt.Errorf("AuthService: Login: jwt is nil")
	}

	return &impl{
		AuthConfig:  cfg,
		AuthRepo:    au,
		AccountRepo: acc,
		JWT:         jwt,
	}, nil
}

func (s *impl) IsAccessValid(ctx context.Context, token string) (bool, error) {
	accountID, err := s.ValidateAccess(ctx, token)
	if err != nil {
		return false, fmt.Errorf("AuthService: ValidateAccess: %w", err)
	}
	return accountID != "", nil
}

func (s *impl) ValidateAccess(ctx context.Context, token string) (string, error) {
	accountID, err := s.JWT.Validate(token)
	if err != nil {
		return "", fmt.Errorf("AuthService: ValidateAccess: %w", err)
	}

	current, err := s.AuthRepo.GetAccess(ctx, accountID)
	if err != nil {
		return "", fmt.Errorf("AuthService: ValidateAccess: %w", err)
	}

	fmt.Println(current)

	if current != token {
		return "", fmt.Errorf("AuthService: ValidateAccess: your token is invalid")
	}

	return accountID, nil
}

func (s *impl) Login(ctx context.Context, accountID string) (*dto.GetTokensDTO, error) {
	if accountID == "" {
		return nil, fmt.Errorf("AuthService: Login: accountId is empty")
	}

	access, refresh, err := s.genTokens(accountID)
	if err != nil {
		return nil, fmt.Errorf("AuthService: Login: %w", err)
	}

	h := sha256.Sum256([]byte(refresh))

	pk, err := s.AccountRepo.GetAccountPK(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("AuthService: Login: %w", err)
	}

	if err := s.AuthRepo.UpsertToken(ctx, models.SaveInput{
		AccountPK:        pk,
		AmoID:            accountID,
		Access:           access,
		RefreshHash:      hex.EncodeToString(h[:]),
		AccessDuration:   s.AuthConfig.AccessTTL,
		RefreshExpiredAt: time.Now().Add(s.AuthConfig.RefreshTTL),
	}); err != nil {
		return nil, fmt.Errorf("AuthService: Login: %w", err)
	}

	return &dto.GetTokensDTO{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (s *impl) RefreshToken(ctx context.Context, token string) (*dto.GetTokensDTO, error) {
	if token == "" {
		return nil, fmt.Errorf("AuthService: RefreshToken: token is empty")
	}

	old, err := s.AuthRepo.GetRefresh(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("AuthService: RefreshToken: %w", err)
	}

	sum := sha256.Sum256([]byte(token))
	newHash := hex.EncodeToString(sum[:])

	if old != newHash {
		return nil, fmt.Errorf("AuthService: RefreshToken: token is invalid")
	}

	// если все норм, то идем создавать
	accountID, err := s.JWT.Validate(token)
	if err != nil {
		return nil, fmt.Errorf("AuthService: RefreshToken: %w", err)
	}

	access, refresh, err := s.genTokens(accountID)
	if err != nil {
		return nil, fmt.Errorf("AuthService: RefreshToken: %w", err)
	}

	h := sha256.Sum256([]byte(refresh))

	// токен истёк -> обновляем
	if err := s.AuthRepo.Update(ctx, models.UpdateInput{
		AccountID:        accountID,
		Access:           access,
		RefreshHash:      hex.EncodeToString(h[:]),
		AccessDuration:   s.AuthConfig.AccessTTL,
		RefreshExpiredAt: time.Now().Add(s.AuthConfig.RefreshTTL),
	}); err != nil {
		return nil, fmt.Errorf("AuthService: RefreshToken: %w", err)
	}

	return &dto.GetTokensDTO{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (s *impl) genTokens(accountID string) (string, string, error) {
	access, err := s.JWT.Generate("access", accountID)
	if err != nil {
		return "", "", err
	}

	refresh, err := s.JWT.Generate("refresh", accountID)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
