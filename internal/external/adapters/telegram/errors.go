package telegram

import (
	"errors"
)

var (
	ErrAdapterInit = errors.New("TelegramAdapter: config is nil")
	// нужно выводить
	ErrTelegramAPI = errors.New("TelegramAdapter:")
)
