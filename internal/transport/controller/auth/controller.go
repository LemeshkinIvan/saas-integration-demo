package auth

import (
	dto "daos_core/internal/domain/dto/auth"
	"daos_core/internal/domain/dto/common"
	"daos_core/internal/domain/services/auth"
	error_handler "daos_core/internal/utils/error"

	"github.com/gin-gonic/gin"
)

// https://www.amocrm.ru/developers/content/oauth/step-by-step - документация, необходимая для понимания контекста
type Controller interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type impl struct {
	service auth.Service
}

func NewAuthCtrl(s auth.Service) Controller {
	return &impl{service: s}
}

func (ctrl *impl) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.LoginDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	data, err := ctrl.service.Login(ctx, request.AccountID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(201, common.RegularResponseDTO[dto.GetTokensDTO]{
		Ok:          true,
		Description: "Token was updated",
		Data:        *data,
	})

}

func (ctrl *impl) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.RefreshTokenDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	data, err := ctrl.service.RefreshToken(ctx, request.Token)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(201, common.RegularResponseDTO[dto.GetTokensDTO]{
		Ok:          true,
		Description: "Token was updated",
		Data:        *data,
	})
}
