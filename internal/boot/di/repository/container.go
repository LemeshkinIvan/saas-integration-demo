package repository

import (
	account_repo "daos_core/internal/data/repositories/account"
	amo_auth "daos_core/internal/data/repositories/amo_tokens"
	"daos_core/internal/data/repositories/auth"
	"daos_core/internal/data/repositories/bot"
	"daos_core/internal/infrastructure/db"

	"daos_core/internal/data/repositories/chat"
	instance_repo "daos_core/internal/data/repositories/instance"
	oauth_repository "daos_core/internal/data/repositories/oauth"
	"daos_core/internal/data/repositories/telegram"
	tg_repo "daos_core/internal/data/repositories/telegram"
)

type Container struct {
	Oauth     oauth_repository.Repository
	Telegram  telegram.Repository
	Instance  instance_repo.Repository
	Account   account_repo.Repository
	Chat      chat.Repository
	Bot       bot.Repository
	AmoTokens amo_auth.Repository
	Auth      auth.Repository
}

func RegisterAll(postgres *db.Postgres, cache *db.RedisStorage) (*Container, error) {
	if postgres == nil {
		return nil, ErrRepositoryInit
	}

	oauth, err := oauth_repository.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	tg, err := tg_repo.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	account, err := account_repo.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	instance, err := instance_repo.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	chat, err := chat.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	botToken, err := bot.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	amoToken, err := amo_auth.NewRepository(postgres)
	if err != nil {
		return nil, err
	}

	auth, err := auth.NewRepository(postgres, cache)
	if err != nil {
		return nil, err
	}

	return &Container{
		Oauth:     oauth,
		Telegram:  tg,
		Account:   account,
		Instance:  instance,
		Chat:      chat,
		Bot:       botToken,
		AmoTokens: amoToken,
		Auth:      auth,
	}, nil
}
