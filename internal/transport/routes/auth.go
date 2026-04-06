package routes

import (
	"daos_core/internal/transport/controller/auth"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	gin  *gin.Engine
	ctrl auth.Controller
	mid  *middleware.Middleware
}

func NewAuthRoutes(g *gin.Engine, ctrl auth.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("AuthRoutes: New: gin is nil")
	}
	return &AuthRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}

func (r *AuthRoutes) RegRoutes() {
	group := r.gin.Group("/auth")

	group.GET("/login", r.ctrl.Login)
	group.POST("/refresh-token", r.ctrl.RefreshToken)
}
