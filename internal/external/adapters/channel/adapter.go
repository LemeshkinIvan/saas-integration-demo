package channel

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	config "daos_core/internal/constants"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const rfc2822Layout = "Mon, 03 Oct 2020 15:11:21 -0700"

type AmoRequestBody struct {
	AccountID      string `json:"account_id"`
	HookApiVersion string `json:"hook_api_version"`
}

type Adapter interface {
	CreateChannel(amojoID string) (*string, error)
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

func (a *impl) CreateChannel(amojoID string) (*string, error) {
	// Content-MD5
	body := AmoRequestBody{
		AccountID:      amojoID,
		HookApiVersion: "v2",
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: CreateChannel: %w", ErrBase, err)
	}

	bodyHash := convertToMd5Hash(jsonBytes)

	// Content Type
	req, err := http.NewRequest("POST", "", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("%w: CreateChannel: %w", ErrBase, err)
	}

	req.Header.Add("Content-MD5", bodyHash)
	//bearer := fmt.Sprintf("Bearer %s", accessToken)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "amoCRM-Chats-Doc-Example/1.0")
	//req.Header.Add("Authorization", bearer)

	// Date
	now := time.Now()
	rfcFormatNow := now.Format(rfc2822Layout)
	req.Header.Add("Date", rfcFormatNow)

	// X - signature (secret)
	apiMethod := fmt.Sprintf(a.constants.ApiConnectMethod, a.constants.ChannelID)
	signature := createSignatureMsg("POST", bodyHash, rfcFormatNow, apiMethod)
	secret := convertToHMACSHA1Hex(signature, a.constants.ChannelSecretKey)
	req.Header.Add("X-Signature", secret)

	// request
	urlSprint := fmt.Sprintf("https://amojo.amocrm.ru/v2/origin/custom/%s/connect", a.constants.ChannelID)
	req.URL, _ = url.Parse(urlSprint)

	response, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: CreateChannel: %w", ErrBase, err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		data := struct {
			ScopeId string `json:"scope_id"`
		}{}

		if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("%w: CreateChannel: %w", ErrBase, err)
		}

		return &data.ScopeId, nil

	case http.StatusUnauthorized:
		return nil, ErrSignature

	case http.StatusBadRequest:
		return nil, ErrChannelExist

	case http.StatusNotFound:
		return nil, ErrUncorrectData

	default:
		return nil, fmt.Errorf("%w: status %d", ErrBase, response.StatusCode)
	}

	// return good response

}

func createSignatureMsg(method string, bodyHash string, rfcFormatNow string, apiMethod string) string {
	parts := []string{
		// POST | GET
		strings.ToUpper(method),
		// checksum
		bodyHash,
		"application/json",
	}

	if rfcFormatNow != "" {
		parts = append(parts, rfcFormatNow)
	}

	parts = append(parts, apiMethod)
	return strings.Join(parts, "\n")
}

func convertToHMACSHA1Hex(message, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func convertToMd5Hash(body []byte) string {
	// считаем MD5
	sum := md5.Sum(body)
	result := hex.EncodeToString(sum[:])
	return result
}
