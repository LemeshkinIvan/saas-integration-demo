package telegram

import (
	"bytes"
	config "daos_core/internal/constants"
	tg_models "daos_core/internal/data/models/api/tg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Adapter interface {
	GetBotInfo(botToken string) ([]byte, error)
	DeleteWebhook(botToken string) (*tg_models.TypicalResponseDto, error)
	GetWebhookInfo(botToken string) ([]byte, error)
	SetWebhook(botToken string) (*tg_models.TypicalResponseDto, error)
}

type impl struct {
	Constants  *config.TelegramConfig
	HttpClient http.Client
}

func NewAdapter(c *config.TelegramConfig) (Adapter, error) {
	if c == nil {
		return nil, ErrAdapterInit
	}
	return &impl{
		Constants:  c,
		HttpClient: http.Client{},
	}, nil
}

func (a *impl) GetBotInfo(botToken string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/getMe", a.Constants, botToken)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := a.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrTelegramAPI, response.StatusCode, string(body))
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a *impl) DeleteWebhook(botToken string) (*tg_models.TypicalResponseDto, error) {
	url := fmt.Sprintf("%s/%s/getMe", a.Constants.BaseURL, botToken)

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := a.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrTelegramAPI, response.StatusCode, string(body))
	}

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var tgResponse tg_models.TypicalResponseDto
	if err := json.Unmarshal(payload, &tgResponse); err != nil {
		return nil, err
	}

	return &tgResponse, nil
}

func (a *impl) GetWebhookInfo(botToken string) ([]byte, error) {
	requestUrl := fmt.Sprintf(
		"%s/%s/getWebhookInfo", a.Constants.BaseURL, botToken)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return data, nil
}

func (a *impl) SetWebhook(botToken string) (*tg_models.TypicalResponseDto, error) {
	requestUrl := fmt.Sprintf("%s%s/setWebhook", a.Constants.BaseURL, botToken)

	bodyBytes, err := json.Marshal(map[string]string{
		"Url": fmt.Sprintf("%s%s", a.Constants.CallbackRoute, botToken),
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestUrl, io.NopCloser(bytes.NewReader(bodyBytes)))
	if err != nil {
		return nil, err
	}

	response, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var tgResponse tg_models.TypicalResponseDto
	if err := json.Unmarshal(payload, &tgResponse); err != nil {
		return nil, err
	}
	return &tgResponse, nil
}
