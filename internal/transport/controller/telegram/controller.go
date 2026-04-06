//https://www.amocrm.ru/developers/content/chats/chat-api-reference - php examples

package telegram

import (
	tg_models "daos_core/internal/data/models/api/tg"
	tg_webhook_models "daos_core/internal/data/models/api/tg/webhooks"
	telegram_request "daos_core/internal/data/models/telegram/request"
	"daos_core/internal/domain/dto/common"
	"daos_core/internal/domain/services/telegram"
	error_handler "daos_core/internal/utils/error"
	custom_logger "daos_core/internal/utils/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	SetWebhook(c *gin.Context)
	ListenWebhook(c *gin.Context) // msg
	GetBotInfo(c *gin.Context)
	GetWebhookInfo(c *gin.Context)
	DeleteWebhook(c *gin.Context)
}

type impl struct {
	service    telegram.Service
	httpClinet *http.Client
}

func NewTgController(s telegram.Service) Controller {
	return &impl{service: s, httpClinet: &http.Client{}}
}

// позволяет получать сообщения автоматом
// нужно указать url на наш сервак для получения msg
// и указать токен бота
// https://tlgrm.ru/docs/bots/api#setwebhook
// POST
func (ctrl *impl) SetWebhook(c *gin.Context) {
	ctx := c.Request.Context()
	var request telegram_request.SendBotTokenDto

	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	data, err := ctrl.service.SetWebhook(ctx, request)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	if data == nil {
		error_handler.HttpErrHandler(c, fmt.Errorf("telegram response is nil"))
		return
	}

	c.JSON(200, common.RegularResponseDTO[map[string]tg_models.TypicalResponseDto]{
		Ok:          data.Ok,
		Description: data.Description,
		Data:        map[string]tg_models.TypicalResponseDto{"telegramResponce": *data},
	})
}

// get last bot activity (message for example)
func (ctrl *impl) ListenWebhook(c *gin.Context) {
	var data tg_webhook_models.GetUpdateDto

	if err := c.ShouldBindJSON(&data); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	custom_logger.AsyncLog(4, string(payload))
	c.Status(200)
}

// такая же как setWebhook, но в url нужно слать пустую строку
// чтобы удалить. либо вообще никакого тела
// POST
func (ctrl *impl) DeleteWebhook(c *gin.Context) {
	ctx := c.Request.Context()
	var request telegram_request.SendBotTokenDto

	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	// se -> adap -> se -> ctrl
	data, err := ctrl.service.DeleteWebhook(ctx, request)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	if data == nil {
		error_handler.HttpErrHandler(c, fmt.Errorf("telegram response is nil"))
		return
	}

	c.JSON(200, common.RegularResponseDTO[map[string]tg_models.TypicalResponseDto]{
		Ok:          data.Ok,
		Description: data.Description,
		Data:        map[string]tg_models.TypicalResponseDto{"telegramResponse": *data},
	})
}

// POST
func (ctrl *impl) GetWebhookInfo(c *gin.Context) {
	// получаем id бота и находим токен у нас
	// отправляем запрос в тг и его парсим
	// возвращаем на фронт
	ctx := c.Request.Context()
	request := struct {
		BotToken string `json:"botToken"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}
	data, err := ctrl.service.GetWebhookInfo(ctx, request.BotToken)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	if data == nil {
		error_handler.HttpErrHandler(c, fmt.Errorf("telegram response is nil"))
		return
	}

	c.JSON(200, common.RegularResponseDTO[[]byte]{
		Ok:          true,
		Description: "",
		Data:        data,
	})
}

// GET (getMe)
func (ctrl *impl) GetBotInfo(c *gin.Context) {
	// получаем id бота и находим токен у нас
	// отправляем запрос в тг и его парсим
	// возвращаем на фронт
	ctx := c.Request.Context()
	request := struct {
		BotToken string `json:"botToken"`
	}{}

	if err := c.ShouldBindJSON(&request); err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	data, err := ctrl.service.GetBotInfo(ctx, request.BotToken)
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
