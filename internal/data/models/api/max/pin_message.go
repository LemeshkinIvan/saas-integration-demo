package max_models

type PinMessageRequest struct {
	MessageId string `json:"message_id"`
	Notify    *bool  `json:"notify"`
}

type PinMessageReponse struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}
