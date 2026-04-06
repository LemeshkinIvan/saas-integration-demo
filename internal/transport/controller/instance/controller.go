package instance

import (
	"daos_core/internal/domain/dto/common"
	dto "daos_core/internal/domain/dto/instance"
	"daos_core/internal/domain/services/instance"
	d "daos_core/internal/transport/controller/instance/dto"
	error_handler "daos_core/internal/utils/error"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetList(c *gin.Context)
	GetByID(c *gin.Context)
	GetCount(c *gin.Context)
	Create(c *gin.Context)
	UpdateByID(c *gin.Context)
	DeleteByID(c *gin.Context)
}

type impl struct {
	service instance.Service
}

func NewInstanceController(s instance.Service) Controller {
	return &impl{service: s}
}

func (ctrl *impl) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	var uri d.GetByIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	accountID := c.GetString("accountId")

	data, err := ctrl.service.GetByID(ctx, accountID, uri.ID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.RegularResponseDTO[dto.GetByIDDTO]{
		Ok:          true,
		Description: "Here you are",
		Data:        *data,
	})
}

func (ctrl *impl) GetList(c *gin.Context) {
	ctx := c.Request.Context()

	accountID := c.GetString("accountId")

	data, err := ctrl.service.ListByAccountID(ctx, accountID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.RegularResponseDTO[dto.GetListDTO]{
		Ok:          true,
		Description: "",
		Data:        *data,
	})
}

func (ctrl *impl) GetCount(c *gin.Context) {
	ctx := c.Request.Context()

	accountID := c.GetString("accountId")

	data, err := ctrl.service.CountInstance(ctx, accountID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.RegularResponseDTO[dto.GetCountDTO]{
		Ok:          true,
		Description: "Big num",
		Data:        *data,
	})
}

func (ctrl *impl) Create(c *gin.Context) {
	ctx := c.Request.Context()

	accountID := c.GetString("accountId")

	data, err := ctrl.service.Create(ctx, accountID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.RegularResponseDTO[dto.GetInstanceDTO]{
		Ok:          true,
		Description: "Instance created successfully",
		Data:        *data,
	})
}

func (ctrl *impl) UpdateByID(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.UpdateDTO
	if err := c.ShouldBindJSON(&data); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	accountID := c.GetString("accountId")

	err := ctrl.service.Update(ctx, data, accountID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.RegularResponseDTO[any]{
		Ok:          true,
		Description: fmt.Sprintf("instance %d was updated", data.InstanceID),
	})
}

func (ctrl *impl) DeleteByID(c *gin.Context) {
	ctx := c.Request.Context()

	var uri d.GetByIDUri
	if err := c.ShouldBindUri(&uri); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	accountID := c.GetString("accountId")

	err := ctrl.service.Delete(ctx, accountID, uri.ID)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(http.StatusNoContent, common.RegularResponseDTO[any]{
		Ok:          true,
		Description: fmt.Sprintf("instance with id = %d was deleted", uri.ID),
	})
}
