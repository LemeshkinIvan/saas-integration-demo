package routes

import (
	"daos_core/internal/transport/controller/account"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AccountRoutes struct {
	gin  *gin.Engine
	ctrl account.Controller
	mid  *middleware.Middleware
}

func NewAccountRoutes(g *gin.Engine, ctrl account.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("AccountRoutes: New: gin is nil")
	}
	return &AccountRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}

func (r *AccountRoutes) RegRoutes() {
	group := r.gin.Group("/accounts")

	group.Use(r.mid.AuthMiddleware())

	group.GET("/", r.ctrl.GetAccount)
	group.GET("/limit", r.ctrl.GetInstanceLimit)
}
