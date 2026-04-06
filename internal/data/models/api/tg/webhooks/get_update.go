package tg_webhook_models

import tg_models "daos_core/internal/data/models/api/tg"

// шлет мне тг. я сократил опциональные поля
type GetUpdateDto struct {
	UpdateId int64              `json:"update_id"`
	Message  *tg_models.Message `json:"message"`
}
