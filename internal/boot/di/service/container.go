package service

import (
	"daos_core/internal/boot/di/adapter"
	"daos_core/internal/boot/di/repository"
	"daos_core/internal/boot/di/utils"
	"daos_core/internal/constants"
	"daos_core/internal/domain/services/account"
	"daos_core/internal/domain/services/auth"
	"daos_core/internal/domain/services/chat"
	"daos_core/internal/domain/services/instance"
	"daos_core/internal/domain/services/oauth"
	"daos_core/internal/domain/services/pipeline"
	"daos_core/internal/domain/services/telegram"
	"daos_core/internal/domain/services/user"

	"fmt"
)

type Container struct {
	Telegram telegram.Service
	Oauth    oauth.Service
	Account  account.Service
	Instance instance.Service
	Pipeline pipeline.Service
	User     user.Service
	Chat     chat.Service
	Auth     auth.Service
}

func RegisterAll(
	r *repository.Container,
	a *adapter.Container,
	cfg *constants.ServiceConfig,
	utils *utils.Container,
) (*Container, error) {

	if err := validateArgs(r, a, cfg, utils); err != nil {
		return nil, err
	}

	auth, err := auth.NewService(&cfg.Auth, r.Auth, r.Account, utils.JWT)
	if err != nil {
		return nil, fmt.Errorf("ServiceContainer: InitService; %w", err)
	}

	return &Container{
		Instance: instance.NewService(r.Instance, r.Bot, r.Account,
			r.AmoTokens, a.Source),

		Account:  account.NewService(r.Account, r.Oauth),
		Oauth:    oauth.NewOauthService(r.Oauth, r.Account, a.Oauth, a.Account, a.Channel),
		Telegram: telegram.NewService(r.Telegram, r.Instance, r.Bot, a.Telegram),
		Pipeline: pipeline.NewImpl(r.AmoTokens, a.Pipeline),
		User:     user.NewImpl(r.AmoTokens, a.Users),
		Chat:     chat.NewChatService(r.Account, a.Chat, r.AmoTokens),
		Auth:     auth,
	}, nil
}

func validateArgs(
	r *repository.Container,
	a *adapter.Container,
	cfg *constants.ServiceConfig,
	utils *utils.Container,
) error {
	if a == nil {
		return fmt.Errorf("%w - adapter container is nil", ErrServiceInit)
	}

	if r == nil {
		return fmt.Errorf("%w - repository container is nil", ErrServiceInit)
	}

	if cfg == nil {
		return fmt.Errorf("%w - cfg container is nil", ErrServiceInit)
	}

	if utils == nil {
		return fmt.Errorf("%w - utils container is nil", ErrServiceInit)
	}

	return nil
}
