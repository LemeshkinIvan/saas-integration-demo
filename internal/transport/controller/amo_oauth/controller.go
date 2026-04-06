package amo_oauth

import (
	"daos_core/internal/domain/dto/common"
	dto "daos_core/internal/domain/dto/oauth"
	"daos_core/internal/domain/services/oauth"
	amo_oauth "daos_core/internal/transport/controller/amo_oauth/dto"
	error_handler "daos_core/internal/utils/error"

	"github.com/gin-gonic/gin"
)

// https://www.amocrm.ru/developers/content/oauth/step-by-step - документация, необходимая для понимания контекста
type Controller interface {
	SaveAmoTokens(c *gin.Context)
	UpdateAmoTokens(c *gin.Context)
}

type impl struct {
	service oauth.Service
}

func NewOauthCtrl(s oauth.Service) Controller {
	return &impl{service: s}
}

func (ctrl *impl) SaveAmoTokens(c *gin.Context) {
	ctx := c.Request.Context()

	var params amo_oauth.SaveAmoTokensQuery
	if err := c.ShouldBindQuery(&params); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	err := ctrl.service.SaveTokens(ctx, params.Code, params.Referer)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(201, common.RegularResponseDTO[any]{
		Ok:          true,
		Description: "Integration connection was successful",
	})
}

func (ctrl *impl) UpdateAmoTokens(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.RefreshTokensDTO
	if err := c.ShouldBindJSON(&data); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	err := ctrl.service.UpdateAccessToken(ctx, data)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(201, common.RegularResponseDTO[any]{
		Ok:          true,
		Description: "Token was updated",
	})
}
