package user

import (
	config "daos_core/internal/constants"
	adapter_http_utils "daos_core/internal/external/common"
	user_dto "daos_core/internal/external/models/user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Adapter interface {
	GetUsers(amoToken string) (*user_dto.GetUsers, error)
}

type impl struct {
	httpClient http.Client
	constants  *config.AmoConfig
}

func NewAdapter(c *config.AmoConfig) (Adapter, error) {
	if c == nil {
		return nil, ErrConfigArgument
	}

	return &impl{
		constants:  c,
		httpClient: http.Client{},
	}, nil
}

func (a *impl) GetUsers(apiToken string) (*user_dto.GetUsers, error) {
	url := fmt.Sprintf("https://%s/api/v4/users", a.constants.Referer)
	//fmt.Println(url)

	request, err := adapter_http_utils.GenDefaultRequest(url, "GET", nil, apiToken)
	if err != nil {
		return nil, err
	}

	// отправляем в amo
	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		// начинаем парсить ответ
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var res user_dto.GetUsers
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return &res, nil

	case http.StatusForbidden:
		return nil, ErrStatusForbidden

	case http.StatusUnauthorized:
		return nil, ErrUnathorized

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}
