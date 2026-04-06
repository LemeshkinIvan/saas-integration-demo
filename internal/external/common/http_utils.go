package adapter_http_utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GenDefaultRequest(url string, method string, body *[]byte, authToken string) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(*body)
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	bearerString := fmt.Sprintf("Bearer %s", authToken)
	request.Header.Add("Authorization", bearerString)

	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

func ConvertResponseBodyToString(response http.Response) (*string, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	// Конвертируем в строку
	bodyString := string(bodyBytes)

	if bodyString == "" {
		bodyString = fmt.Sprintf("there no body msg, folk. status code %d", response.StatusCode)
	}

	return &bodyString, nil
}

func GenAmoXSignature(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, fmt.Sprint(p))
	}

	return strings.Join(out, "\n")
}

// later
// func GenAmoRequestWithSignature(method string, url string, body []byte) (*http.Request, error) {
// 	return nil, nil
// }
