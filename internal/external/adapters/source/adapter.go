package source

import (
	config "daos_core/internal/constants"
	dto "daos_core/internal/domain/dto/instance"
	adapter_http_utils "daos_core/internal/external/common"
	models "daos_core/internal/external/models/source"
	"daos_core/internal/utils/debug"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Adapter interface {
	UpdateByID(auth models.AuthInput, data dto.UpdateDTO) error
	Create(auth models.AuthInput, input models.CreateInput) (*models.GetSourcesResponseDTO, error)
	DeleteByID(auth models.AuthInput, sourceID int) error
}

type impl struct {
	Constants  *config.AmoConfig
	HttpClient http.Client
}

func NewAdapter(c *config.AmoConfig) (Adapter, error) {
	if c == nil {
		return nil, fmt.Errorf("SourceAdapter: config is nil")
	}
	return &impl{
		HttpClient: http.Client{},
		Constants:  c,
	}, nil
}

func (c *impl) UpdateByID(auth models.AuthInput, data dto.UpdateDTO) error {
	url := fmt.Sprintf("https://%s/%s/%d", auth.Referer, c.Constants.SourceURL, data.SourceID)

	body := struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		PipelineID int    `json:"pipeline_id"`
	}{*data.SourceID, data.Name, *data.PipelineID}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("SourceAdapter: UpdateByID: %w", err)
	}

	// создаем запрос
	request, err := adapter_http_utils.GenDefaultRequest(url, "PATCH", &bodyBytes, auth.Referer)
	if err != nil {
		return fmt.Errorf("SourceAdapter: UpdateByID: %w", err)
	}

	// отправляем в amo
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("SourceAdapter: UpdateByID: %w", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return nil

	case http.StatusUnauthorized:
		return fmt.Errorf("SourceAdapter: UpdateByID: %w", ErrUnautharized)

	case http.StatusBadRequest:
		return fmt.Errorf("SourceAdapter: UpdateByID: %w", ErrUncorrectData)

	default:
		return fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}

func (c *impl) Create(auth models.AuthInput, input models.CreateInput) (*models.GetSourcesResponseDTO, error) {
	url := "https://" + auth.Referer + "/" + c.Constants.SourceURL

	// инициализируем dto для body
	dto := genCreateSourceDto(input.Name, input.ExternalID)
	debug.Dump(dto)

	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("SourceAdapter: Create: %w", err)
	}

	// создаем запрос
	request, err := adapter_http_utils.GenDefaultRequest(url, "POST", &dtoBytes, auth.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("SourceAdapter: Create: %w", err)
	}

	// отправляем в amo
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("SourceAdapter: Create: %w", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		var data models.GetSourcesResponseDTO
		if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("SourceAdapter: Create: %w", err)
		}

		return &data, nil

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("SourceAdapter: Create: %w", ErrUnautharized)

	case http.StatusBadRequest:
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("SourceAdapter: Create: %w", err)
		}

		fmt.Println(string(body))
		return nil, fmt.Errorf("SourceAdapter: Create: %w", ErrUncorrectData)

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}

func (c *impl) DeleteByID(auth models.AuthInput, sourceID int) error {
	strID := strconv.Itoa(sourceID)

	url := "https://" + auth.Referer + c.Constants.SourceURL + "/" + "?id=" + strID

	fmt.Println(url)
	request, err := adapter_http_utils.GenDefaultRequest(url, "DELETE", nil, auth.AccessToken)
	if err != nil {
		return fmt.Errorf("SourceAdapter: DeleteByID: %w", err)
	}

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("SourceAdapter: DeleteByID: %w", err)
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil

	case http.StatusUnauthorized:
		return fmt.Errorf("SourceAdapter: DeleteByID: %w", ErrUnautharized)

	case http.StatusBadRequest:
		return fmt.Errorf("SourceAdapter: DeleteByID: %w", ErrUncorrectData)

	case http.StatusNotFound:
		return fmt.Errorf("SourceAdapter: DeleteByID: %w", ErrNotFound)

	default:
		f, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("SourceAdapter: DeleteByID: %w", err)
		}

		fmt.Println(string(f[:]))
		return fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}

func genCreateSourceDto(name string, externalId string) models.CreateSourceRequestDTO {
	source := models.SourceRequest{
		Name:       name,
		ExternalID: externalId,
	}

	return models.CreateSourceRequestDTO{Data: source}
}
