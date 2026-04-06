package instance

type UpdateDTO struct {
	InstanceID int64        `json:"instanceId" binding:"required"`
	Name       string       `json:"name" binding:"required"`
	Bot        UpdateBotDTO `json:"bot" binding:"required"`
	PipelineID *int         `json:"pipelineId" binding:"required"`
	SourceID   *int         `json:"sourceId" binding:"required"`
}

type UpdateBotDTO struct {
	Type  int8   `json:"type" binding:"required"`
	Token string `json:"token" binding:"required"`
}
