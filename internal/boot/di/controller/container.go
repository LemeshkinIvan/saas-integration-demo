package controller

import (
	"daos_core/internal/boot/di/service"
	"daos_core/internal/transport/controller/account"
	"daos_core/internal/transport/controller/amo_oauth"
	"daos_core/internal/transport/controller/auth"
	"daos_core/internal/transport/controller/chat"
	"daos_core/internal/transport/controller/instance"
	"daos_core/internal/transport/controller/pipeline"
	"daos_core/internal/transport/controller/telegram"
	"fmt"
)

type Container struct {
	Oauth    amo_oauth.Controller
	Telegram telegram.Controller
	Account  account.Controller
	Instance instance.Controller
	Pipeline pipeline.Controller
	Chat     chat.Controller
	Auth     auth.Controller
}

func RegisterAll(s *service.Container) (*Container, error) {
	if s == nil {
		return nil, fmt.Errorf("%w - service is nil", ErrControllerInit)
	}

	return &Container{
		Oauth:    amo_oauth.NewOauthCtrl(s.Oauth),
		Telegram: telegram.NewTgController(s.Telegram),
		Instance: instance.NewInstanceController(s.Instance),
		Account:  account.NewAccountController(s.Account),
		Pipeline: pipeline.NewPipelineController(s.Pipeline),
		Chat:     chat.NewChatController(s.Chat),
		Auth:     auth.NewAuthCtrl(s.Auth),
	}, nil
}
