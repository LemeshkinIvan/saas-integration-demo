package user

import (
	"context"
	amo_auth "daos_core/internal/data/repositories/amo_tokens"

	user_adapter "daos_core/internal/external/adapters/users"
	user_dto "daos_core/internal/external/models/user"
	"fmt"
)

type Service interface {
	GetUsers(ctx context.Context, accountID string) (*user_dto.GetUsers, error)
}

type impl struct {
	AmoTokensRepo amo_auth.Repository
	UsersAdapter  user_adapter.Adapter
}

func NewImpl(r amo_auth.Repository, a user_adapter.Adapter) Service {
	return &impl{
		AmoTokensRepo: r,
		UsersAdapter:  a,
	}
}

func (s *impl) GetUsers(ctx context.Context, accountID string) (*user_dto.GetUsers, error) {
	token, err := s.AmoTokensRepo.GetAccessToken(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if token == "" {
		return nil, fmt.Errorf("amo token not found")
	}

	data, err := s.UsersAdapter.GetUsers(token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no users received from amo")
	}
	return data, err
}
