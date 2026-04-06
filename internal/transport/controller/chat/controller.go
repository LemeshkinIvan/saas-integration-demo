package chat

import (
	dto "daos_core/internal/domain/dto/chat"
	"daos_core/internal/domain/dto/common"
	"daos_core/internal/domain/services/chat"
	error_handler "daos_core/internal/utils/error"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Create(c *gin.Context)
}

type impl struct {
	service chat.Service
}

func NewChatController(s chat.Service) Controller {
	return &impl{service: s}
}

func (ctrl *impl) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.CreateChatDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	err := ctrl.service.Create(ctx, request)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(200, common.RegularResponseDTO[any]{
		Ok:          true,
		Description: "chat was created",
	})
}
