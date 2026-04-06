package adapter

import (
	config "daos_core/internal/constants"
	"daos_core/internal/external/adapters/account"
	"daos_core/internal/external/adapters/channel"
	chat_adapter "daos_core/internal/external/adapters/chat"
	oauth_adapter "daos_core/internal/external/adapters/oauth"
	pipeline_adapter "daos_core/internal/external/adapters/pipeline"
	source_adapter "daos_core/internal/external/adapters/source"
	telegram_adapter "daos_core/internal/external/adapters/telegram"
	user_adapter "daos_core/internal/external/adapters/users"
	"fmt"
)

type Container struct {
	Source   source_adapter.Adapter
	Oauth    oauth_adapter.Adapter
	Channel  channel.Adapter
	Account  account.Adapter
	Pipeline pipeline_adapter.Adapter
	Users    user_adapter.Adapter
	Telegram telegram_adapter.Adapter
	Chat     chat_adapter.Adapter
}

// обращение к различным API: telegram, max, amo
func RegisterAll(a_c *config.AmoConfig, t_c *config.TelegramConfig) (*Container, error) {
	if a_c == nil {
		return nil, ErrAmoConfig
	}

	if t_c == nil {
		return nil, ErrTelegramConfig
	}

	source, err := source_adapter.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	oauth, err := oauth_adapter.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	pipeline, err := pipeline_adapter.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	users, err := user_adapter.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	telegram, err := telegram_adapter.NewAdapter(t_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	chat, err := chat_adapter.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	account, err := account.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	channel, err := channel.NewAdapter(a_c)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	return &Container{
		Source:   source,
		Oauth:    oauth,
		Pipeline: pipeline,
		Users:    users,
		Telegram: telegram,
		Account:  account,
		Chat:     chat,
		Channel:  channel,
	}, nil
}

func returnFormatErr(e error) error {
	return fmt.Errorf("%w \n%w", ErrAdaptersInit, e)
}
