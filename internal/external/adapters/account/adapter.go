package account

import (
	config "daos_core/internal/constants"
	auth_models "daos_core/internal/data/models/oauth"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Adapter interface {
	GetAccount(token string, referer string) (*auth_models.AmoAccountOutput, error)
}

type impl struct {
	httpClient http.Client
	constants  *config.AmoConfig
}

func NewAdapter(c *config.AmoConfig) (Adapter, error) {
	if c == nil {
		return nil, ErrAdapterInit
	}

	return &impl{
		constants:  c,
		httpClient: http.Client{},
	}, nil
}

func (a *impl) GetAccount(token string, referer string) (*auth_models.AmoAccountOutput, error) {
	// создание запроса в AMO для получение данных аккаунта
	getUserUrl := fmt.Sprintf("https://%s/api/v4/account?with=amojo_id", referer)
	bearer := fmt.Sprintf("Bearer %s", token)

	request, err := http.NewRequest("GET", getUserUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBase, err)
	}

	request.Header.Add("Authorization", bearer)
	// запрос на получение пользака
	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBase, err)
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var body auth_models.AmoAccount
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBase, err)
		}

		f := strconv.Itoa(body.AmoID)
		return &auth_models.AmoAccountOutput{
			AmoID:     f,
			Name:      body.Name,
			Subdomain: body.Subdomain,
			AmojoID:   body.AmojoID,
		}, nil

	case http.StatusUnauthorized:
		return nil, ErrUnauthorized

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}
