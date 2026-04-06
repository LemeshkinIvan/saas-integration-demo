package oauth

import (
	"bytes"
	config "daos_core/internal/constants"
	auth_models "daos_core/internal/data/models/oauth"
	models "daos_core/internal/external/models/oauth"
	"encoding/json"
	"fmt"
	"net/http"
)

type Adapter interface {
	GetTokens(code string, referer string) (*auth_models.TokensResponse, error)
	RefreshTokens(data models.RefreshInput) (*models.UpdateOutput, error)
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

func (a *impl) RefreshTokens(data models.RefreshInput) (*models.UpdateOutput, error) {
	uri := fmt.Sprintf("https://%s.amocrm.ru/oauth2/access_token", a.constants.Domain)

	body := map[string]string{
		"client_id":     a.constants.ChannelIDOauth,
		"client_secret": a.constants.ChannelSecret,
		"grant_type":    "refresh_token",
		"refresh_token": data.Refresh,
		"redirect_uri":  a.constants.RedirectURL,
	}

	body_bytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: RefreshTokens: %w", ErrBase, err)
	}

	// создание запроса в AMO на получение tokens
	request, err := http.NewRequest("POST", uri, bytes.NewReader(body_bytes))
	if err != nil {
		return nil, fmt.Errorf("%w: RefreshTokens: %w", ErrBase, err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "amoCRM-oAuth-client/1.0")

	// получаем токенс
	resp, err := a.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: RefreshTokens: %w", ErrBase, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var payload models.UpdateTokensResponseDTO
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("%w: RefreshTokens: %w", ErrBase, err)
		}

		// fmt.Printf("%+v\n", payload)

		return &models.UpdateOutput{
			RefreshToken: payload.RefreshToken,
			AccessToken:  payload.AccessToken,
			Referer:      data.Referer,
			ExpiresIn:    payload.ExpiresIn,
		}, nil

	case http.StatusUnauthorized:
		return nil, ErrUnautharized

	case http.StatusBadRequest:
		return nil, ErrUncorrectData

	case http.StatusNotFound:
		return nil, ErrNotFound

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, resp.StatusCode)
	}
}

func (a *impl) GetTokens(code string, referer string) (*auth_models.TokensResponse, error) {
	uri := fmt.Sprintf("https://%s.amocrm.ru/oauth2/access_token", a.constants.Domain)

	body := map[string]string{
		"client_id":     a.constants.ChannelIDOauth,
		"client_secret": a.constants.ChannelSecret,
		"code":          code,
		"redirect_uri":  a.constants.RedirectURL,
		"grant_type":    "authorization_code",
	}

	body_bytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: GetTokens: %w", ErrBase, err)
	}

	// создание запроса в AMO на получение tokens
	request, err := http.NewRequest("POST", uri, bytes.NewReader(body_bytes))
	if err != nil {
		return nil, fmt.Errorf("%w: GetTokens: %w", ErrBase, err)
	}

	request.Header.Add("Content-Type", "application/json")

	// получаем токенс
	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: GetTokens: %w", ErrBase, err)
	}

	// закрываем соединение
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var payload auth_models.TokensResponse
		if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("%w: GetTokens: %w", ErrBase, err)
		}
		return &payload, nil

	case http.StatusBadRequest:
		return nil, ErrIncorrectData

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}
