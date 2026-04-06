package routes

import (
	"daos_core/internal/transport/controller/pipeline"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PipelineRoutes struct {
	gin  *gin.Engine
	ctrl pipeline.Controller
	mid  *middleware.Middleware
}

func NewPipelineRoutes(g *gin.Engine, ctrl pipeline.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("PipelineRoutes: New: gin is nil")
	}
	return &PipelineRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}

func (r *PipelineRoutes) RegRoutes() {
	r.gin.GET("/pipeline", r.mid.AuthMiddleware(), r.ctrl.GetPipeline)
}
