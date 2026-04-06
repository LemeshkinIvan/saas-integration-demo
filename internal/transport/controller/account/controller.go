package account

import (
	"daos_core/internal/domain/dto/account"
	"daos_core/internal/domain/dto/common"
	srv "daos_core/internal/domain/services/account"
	error_handler "daos_core/internal/utils/error"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetAccount(c *gin.Context)
	GetInstanceLimit(c *gin.Context)
	// maybe later
	//UpdateAccount(c *gin.Context)
	GetTokens(c *gin.Context)
	UpdateTokens(c *gin.Context)
}

type impl struct {
	service srv.Service
}

func NewAccountController(s srv.Service) Controller {
	return &impl{service: s}
}

func (ctrl *impl) GetInstanceLimit(c *gin.Context) {
	ctx := c.Request.Context()

	accountId := c.GetString("accountId")

	dto, err := ctrl.service.GetInstanceLimit(ctx, accountId)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(200, common.RegularResponseDTO[account.GetInstanceLimitDTO]{
		Ok:          true,
		Description: "",
		Data:        *dto,
	})
}

func (ctrl *impl) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()

	accountId := c.GetString("accountId")

	account, err := ctrl.service.Get(ctx, accountId)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(200, account)
}

func (ctrl *impl) UpdateTokens(c *gin.Context) {}

func (ctrl *impl) GetTokens(c *gin.Context) {
	ctx := c.Request.Context()

	accountId := c.GetString("accountId")

	pair, err := ctrl.service.GetTokens(ctx, accountId)
	if err != nil {
		error_handler.HttpErrHandler(c, err)
		return
	}

	c.JSON(200, pair)
}

// TODO: refactor this to adapter
// func (ctrl *impl) UpdateAccount(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	data := struct {
// 		Referer string `json:"referer"`
// 	}{}

// 	if err := c.ShouldBindJSON(&data); err != nil {
// 		api_utils.HandleError(400, err.Error(), c)
// 		return
// 	}

// 	amoToken, err := ctrl.service.GetAmoApiAccessToken(ctx, data.Referer)
// 	if err != nil {
// 		api_utils.HandleError(404, err.Error(), c)
// 		return
// 	}

// 	getUserUrl := fmt.Sprintf("https://%s/api/v4/account?with=amojo_id", data.Referer)
// 	bearer := fmt.Sprintf("Bearer %s", *amoToken)

// 	getUserRequest, err := http.NewRequest("GET", getUserUrl, nil)
// 	if err != nil {
// 		api_utils.HandleError(500, err.Error(), c)
// 		return
// 	}

// 	getUserRequest.Header.Add("Authorization", bearer)
// 	// запрос на получение пользака
// 	userResponse, err := ctrl.httpClient.Do(getUserRequest)
// 	if err != nil {
// 		api_utils.HandleError(500, err.Error(), c)
// 		return
// 	}

// 	defer userResponse.Body.Close()

// 	// паршу user
// 	// обычный c.ShouldBindJson не подойдет, ибо отправлял запрос
// 	var account auth_models.AmoAccount
// 	if err := json.NewDecoder(userResponse.Body).Decode(&account); err != nil {
// 		api_utils.HandleError(400, err.Error(), c)
// 		return
// 	}

// 	if err := ctrl.service.UpdateAccount(ctx, data.Referer, account); err != nil {
// 		api_utils.HandleError(400, err.Error(), c)
// 		return
// 	}
// }
