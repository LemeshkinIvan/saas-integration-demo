package routes

import (
	"daos_core/internal/transport/controller/instance"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type InstanceRoutes struct {
	gin  *gin.Engine
	ctrl instance.Controller
	mid  *middleware.Middleware
}

func NewInstanceRoutes(g *gin.Engine, ctrl instance.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("InstanceRoutes: New: gin is nil")
	}
	return &InstanceRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}

func (r *InstanceRoutes) RegRoutes() {
	group := r.gin.Group("/instances")

	group.Use(r.mid.AuthMiddleware())

	group.GET("/", r.ctrl.GetList)
	group.GET("/:id", r.ctrl.GetByID)
	group.GET("/count", r.ctrl.GetCount)
	group.POST("/", r.ctrl.Create)
	group.PUT("/", r.ctrl.UpdateByID)
	group.DELETE("/:id", r.ctrl.DeleteByID)
}
