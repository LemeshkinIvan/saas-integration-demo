package pipeline

import (
	"daos_core/internal/domain/dto/common"
	"daos_core/internal/domain/services/pipeline"
	error_handler "daos_core/internal/utils/error"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetPipeline(c *gin.Context)
}

type impl struct {
	service pipeline.Service
}

func NewPipelineController(s pipeline.Service) Controller {
	return &impl{service: s}
}

func (ctrl *impl) GetPipeline(c *gin.Context) {
	ctx := c.Request.Context()

	accountId := c.GetString("accountId")

	data, err := ctrl.service.Get(ctx, accountId)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(200, common.RegularResponseDTO[[]byte]{
		Ok:          true,
		Description: "",
		Data:        data,
	})
}
