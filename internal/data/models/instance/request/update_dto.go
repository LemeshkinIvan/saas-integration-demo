package instance_request

type UpdateRequestDTO struct {
	InstanceID int    `json:"instanceId"`
	AccountID  int    `json:"accountId"`
	Name       string `json:"name"`
	BotToken   string `json:"botToken"`
	PipelineID int    `json:"pipeline_id"`
	SourceID   int    `json:"source_id"`
}
