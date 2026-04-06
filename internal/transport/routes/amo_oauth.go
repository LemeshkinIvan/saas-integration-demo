package routes

import (
	"daos_core/internal/transport/controller/amo_oauth"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AmoOauthRoutes struct {
	gin  *gin.Engine
	ctrl amo_oauth.Controller
	mid  *middleware.Middleware
}

func NewAmoOauthRoutes(g *gin.Engine, ctrl amo_oauth.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("AmoOauthRoutes: New: gin is nil")
	}
	return &AmoOauthRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}

func (r *AmoOauthRoutes) RegRoutes() {
	group := r.gin.Group("/oauth")

	group.GET("/callback", r.ctrl.SaveAmoTokens)
	group.POST("/updateToken", r.ctrl.UpdateAmoTokens)
}
