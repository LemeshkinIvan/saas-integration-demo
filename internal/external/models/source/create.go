package source

type CreateInput struct {
	Name       string
	ExternalID string
}

type AuthInput struct {
	Referer     string
	AccessToken string
}

// шлет фронт/amo
type CreateSourcesDTO struct {
	Data []SourceRequest `json:"data"`
}

type CreateSourceRequestDTO struct {
	Data SourceRequest `json:"data"`
}

type SourceRequest struct {
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
}

type Param struct {
	Waba                   bool `json:"waba,omitempty"`
	IsSupportedListMessage bool `json:"is_supported_list_message,omitempty"`
}

type Page struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}
