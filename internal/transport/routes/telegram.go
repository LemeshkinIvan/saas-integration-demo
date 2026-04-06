package routes

import (
	"daos_core/internal/transport/controller/telegram"
	"daos_core/internal/transport/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TelegramRoutes struct {
	gin  *gin.Engine
	ctrl telegram.Controller
	mid  *middleware.Middleware
}

func NewTelegramRoutes(g *gin.Engine, ctrl telegram.Controller, mid *middleware.Middleware) (APIRoute, error) {
	if g == nil {
		return nil, fmt.Errorf("TelegramRoutes: New: gin is nil")
	}
	return &TelegramRoutes{
		gin:  g,
		ctrl: ctrl,
		mid:  mid,
	}, nil
}

func (r *TelegramRoutes) RegRoutes() {
	group := r.gin.Group("/telegram")

	//api.engine.POST("/telegram/sendMessage", api.controller.SendMessage)
	group.POST("/getBotInfo", r.ctrl.GetBotInfo)
	// при подключении пользовательского бота
	group.POST("/webhook", r.ctrl.SetWebhook)
	group.POST("/listenWebhook/:botToken", r.ctrl.ListenWebhook)
	group.POST("/getWebhookInfo", r.ctrl.GetWebhookInfo)
	group.DELETE("/webhook", r.ctrl.DeleteWebhook)
}
