package pipeline

import (
	config "daos_core/internal/constants"
	adapter_http_utils "daos_core/internal/external/common"
	"fmt"
	"io"
	"net/http"
)

type Adapter interface {
	GetPipeline(amoToken string) ([]byte, error)
	//CreatePipeline(pipeline)
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
		httpClient: http.Client{},
		constants:  c,
	}, nil
}

func (a *impl) GetPipeline(apiToken string) ([]byte, error) {
	url := fmt.Sprintf("https://%s%s", a.constants.Referer, a.constants.PipelineURL)
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
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return body, nil

	case http.StatusUnauthorized:
		return nil, ErrUnauthorized

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}
