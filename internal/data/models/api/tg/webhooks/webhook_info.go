package tg_webhook_models

// шлет нам тг
type GetWebhookInfo struct {
	Url                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	LastErrorDate        int    `json:"last_error_date,omitempty"`
	LastErrorMessage     string `json:"last_error_message,omitempty"`
}
