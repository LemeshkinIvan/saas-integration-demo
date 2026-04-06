// //https://www.amocrm.ru/developers/content/chats/chat-api-reference - php examples

package max

// import (
// 	//models "amo_back/internal/data/models"

// 	custom_error "daos_core/internal/common/error"
// 	custom_logger "daos_core/internal/common/logger"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // unnecessary
// type ChatController interface {
// 	GetChats(c *gin.Context)
// 	GetChatById(c *gin.Context)
// 	DeleteChatById(c *gin.Context)

// 	GetPinnedMessage(c *gin.Context)
// 	PinMessage(c *gin.Context)
// 	DeletePinMessage(c *gin.Context)
// }

// type chatController struct {
// 	repository repo.MaxRepository
// 	constants  *chatConstants
// 	client     *http.Client
// }

// func NewChatController(repo repo.MaxRepository, c *chatConstants) ChatController {
// 	return &chatController{
// 		repository: repo,
// 		client:     &http.Client{},
// 		constants:  c,
// 	}
// }

// func (ctrl *chatController) GetChats(c *gin.Context) {
// 	token_url := fmt.Sprintf("https://%s/chats/", ctrl.constants.BASE_URL)

// 	req, err := http.NewRequest("GET", token_url, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	q := req.URL.Query()
// 	q.Add("access_token", ctrl.constants.BOT_TOKEN)

// 	req.URL.RawQuery = q.Encode()

// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetChats").Error())
// 		return
// 	}

// 	resp, err := ctrl.client.Do(req)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetChats").Error())
// 		return
// 	}

// 	c.JSON(200, resp)
// }

// func (ctrl *chatController) GetChatById(c *gin.Context) {
// 	chat_id := c.Query("id")
// 	token_url := fmt.Sprintf("https://%s/chats/%s", ctrl.constants.BASE_URL, chat_id)

// 	req, err := http.NewRequest("GET", token_url, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	q := req.URL.Query()
// 	q.Add("access_token", ctrl.constants.BOT_TOKEN)

// 	req.URL.RawQuery = q.Encode()

// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetChatById").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	resp, err := ctrl.client.Do(req)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetChatById").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	c.JSON(200, resp)
// }

// func (ctrl *chatController) DeleteChatById(c *gin.Context) {
// 	chat_id := c.Query("id")

// 	_, err := strconv.ParseInt(chat_id, 10, 64)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "DeleteChatById").Error())
// 		c.JSON(404, err)
// 		return
// 	}

// 	token_url := fmt.Sprintf("https://%s/chats/%s", ctrl.constants.BASE_URL, chat_id)

// 	req, err := http.NewRequest("DELETE", token_url, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	q := req.URL.Query()
// 	q.Add("access_token", ctrl.constants.BOT_TOKEN)

// 	req.URL.RawQuery = q.Encode()

// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "DeleteChatById").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	resp, err := ctrl.client.Do(req)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "DeleteChatById").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	c.JSON(200, resp)
// }

// func (ctrl *chatController) GetPinnedMessage(c *gin.Context) {
// 	chat_id := c.Query("id")

// 	_, err := strconv.ParseInt(chat_id, 10, 64)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetPinnedMessage").Error())
// 		c.JSON(404, err)
// 		return
// 	}

// 	token_url := fmt.Sprintf("https://%s/chats/%s/pin", ctrl.constants.BASE_URL, chat_id)

// 	req, err := http.NewRequest("GET", token_url, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	q := req.URL.Query()
// 	q.Add("access_token", ctrl.constants.BOT_TOKEN)

// 	req.URL.RawQuery = q.Encode()

// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetPinnedMessage").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	resp, err := ctrl.client.Do(req)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "GetPinnedMessage").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	c.JSON(200, resp)
// }

// func (ctrl *chatController) PinMessage(c *gin.Context) {
// 	chat_id := c.Query("id")

// 	_, err := strconv.ParseInt(chat_id, 10, 64)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "PinMessage").Error())
// 		c.JSON(404, err)
// 		return
// 	}

// 	token_url := fmt.Sprintf("https://%s/chats/%s/pin", ctrl.constants.BASE_URL, chat_id)

// 	req, err := http.NewRequest("PUT", token_url, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	q := req.URL.Query()
// 	q.Add("access_token", ctrl.constants.BOT_TOKEN)

// 	req.URL.RawQuery = q.Encode()

// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "PinMessage").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	resp, err := ctrl.client.Do(req)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "PinMessage").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	c.JSON(200, resp)
// }

// func (ctrl *chatController) DeletePinMessage(c *gin.Context) {
// 	chat_id := c.Query("id")

// 	_, err := strconv.ParseInt(chat_id, 10, 64)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "DeletePinMessage").Error())
// 		c.JSON(404, err)
// 		return
// 	}

// 	token_url := fmt.Sprintf("https://%s/chats/%s/pin", ctrl.constants.BASE_URL, chat_id)

// 	req, err := http.NewRequest("DELETE", token_url, nil)
// 	req.Header.Add("Content-Type", "application/json")

// 	q := req.URL.Query()
// 	q.Add("access_token", ctrl.constants.BOT_TOKEN)

// 	req.URL.RawQuery = q.Encode()

// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "DeletePinMessage").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	resp, err := ctrl.client.Do(req)
// 	if err != nil {
// 		custom_logger.Logg.Warn(custom_error.GetError(err, "DeletePinMessage").Error())
// 		c.JSON(500, err)
// 		return
// 	}

// 	c.JSON(200, resp)
// }
