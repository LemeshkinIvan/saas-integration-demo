package tg_webhook_models

import tg_models "daos_core/internal/data/models/api/tg"

type ResponseSetWebhook struct {
	TelegramResponse tg_models.TypicalResponseDto `json:"telegramResponse"`
}
