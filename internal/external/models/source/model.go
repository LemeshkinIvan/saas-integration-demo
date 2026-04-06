package source

type GetSourcesResponseDTO struct {
	Embedded Embedded `json:"_embedded"`
}

type Embedded struct {
	Sources []SourceResponse `json:"sources"`
}

type SourceResponse struct {
	ID         int    `json:"id"`
	PipelineID *int   `json:"pipeline_id"`
	ExternalID string `json:"external_id"`
}
