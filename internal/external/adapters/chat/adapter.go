package chat

import (
	"crypto/md5"
	config "daos_core/internal/constants"

	adapter_http_utils "daos_core/internal/external/common"
	model "daos_core/internal/external/models/chat"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Adapter interface {
	Connect() error
	Disconnect() error
	Create(input model.CreateInput) error
	Send()
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

func (a *impl) Connect() error {
	return nil
}

func (a *impl) Disconnect() error {
	return nil
}

func (a *impl) Create(input model.CreateInput) error {
	path := fmt.Sprintf("/v2/origin/custom/%s/chats", input.ScopeID)
	url := fmt.Sprintf("https://%s/%s", input.Referer, path)

	body := model.CreateChatAmo{
		ConversationID: input.ConversationID,
		User: model.User{
			ID: input.UserID,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := adapter_http_utils.GenDefaultRequest(url, "POST", &bodyBytes, input.Token)
	if err != nil {
		return err
	}

	date := time.Now().Format(time.RFC1123Z)
	checksum := md5.Sum(bodyBytes)
	hash := checksum[:]

	signature := adapter_http_utils.GenAmoXSignature("POST", string(hash), "application/json", date, path)
	request.Header.Add("X-Signature", signature)
	request.Header.Add("Content-MD5", string(hash))
	request.Header.Add("Date", date)

	response, err := a.httpClient.Do(request)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return nil

	case http.StatusForbidden:
		return ErrSignatureInvalid

	case http.StatusBadRequest:
		return ErrUncorrectData

	default:
		return fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}
}

func (a *impl) Send() {}
