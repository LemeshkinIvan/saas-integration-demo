package max_models

type Icon struct {
	Url string `json:"url"`
}

type PhotoAttachmentPayload struct {
	PhotoId int64  `json:"photo_id"`
	Token   string `json:"token"`
	Url     string `json:"url"`
}
