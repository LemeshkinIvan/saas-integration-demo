package route

import (
	"fmt"

	"daos_core/internal/boot/di/controller"
	"daos_core/internal/transport/middleware"
	"daos_core/internal/transport/routes"

	"github.com/gin-gonic/gin"
)

type Container struct {
	Telegram routes.APIRoute
	Account  routes.APIRoute // required
	Instance routes.APIRoute // required
	Oauth    routes.APIRoute // required
	Pipeline routes.APIRoute // required
	Chat     routes.APIRoute // required
	Auth     routes.APIRoute
}

func RegisterAll(
	g *gin.Engine,
	ctrl *controller.Container,
	mid *middleware.Middleware,
) (*Container, error) {
	if err := validateArgs(g, ctrl, mid); err != nil {
		return nil, err
	}

	telegram, err := routes.NewTelegramRoutes(g, ctrl.Telegram, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	pipeline, err := routes.NewPipelineRoutes(g, ctrl.Pipeline, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	account, err := routes.NewAccountRoutes(g, ctrl.Account, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	instance, err := routes.NewInstanceRoutes(g, ctrl.Instance, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	oauth, err := routes.NewAmoOauthRoutes(g, ctrl.Oauth, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	chat, err := routes.NewChatRoutes(g, ctrl.Chat, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	auth, err := routes.NewAuthRoutes(g, ctrl.Auth, mid)
	if err != nil {
		return nil, returnFormatErr(err)
	}

	telegram.RegRoutes()
	auth.RegRoutes()
	chat.RegRoutes()
	oauth.RegRoutes()
	instance.RegRoutes()
	account.RegRoutes()
	pipeline.RegRoutes()

	return &Container{
		Telegram: telegram,
		Pipeline: pipeline,
		Account:  account,
		Instance: instance,
		Oauth:    oauth,
		Chat:     chat,
		Auth:     auth,
	}, nil
}

func validateArgs(
	g *gin.Engine,
	ctrl *controller.Container,
	mid *middleware.Middleware,
) error {
	if g == nil {
		return fmt.Errorf("%w: gin is nil", ErrInvalidArgument)
	}

	if ctrl == nil {
		return fmt.Errorf("%w: controller container is nil", ErrInvalidArgument)
	}

	if mid == nil {
		return fmt.Errorf("%w: middleware  is nil", ErrInvalidArgument)
	}
	return nil
}

func returnFormatErr(e error) error {
	return fmt.Errorf("%w \n%w", ErrAPIInit, e)
}
