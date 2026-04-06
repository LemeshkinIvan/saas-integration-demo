package routes

import (
	"daos_core/internal/transport/controller/chat"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ChatRoutes struct {
	gin  *gin.Engine
	ctrl chat.Controller
	mid  *middleware.Middleware
}

func NewChatRoutes(g *gin.Engine, ctrl chat.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("ChatRoutes: New: gin is nil")
	}
	return &ChatRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}
func (r *ChatRoutes) RegRoutes() {
	r.gin.POST("/chat", r.mid.AuthMiddleware(), r.ctrl.Create)
}
