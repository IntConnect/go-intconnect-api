package model

type CreatePipelineDto struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Config      map[string]interface{} `json:"config"`
}

type UpdatePipelineDto struct {
	Id          uint64                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Config      map[string]interface{} `json:"config"`
}

type PipelineResponse struct {
	ID                 uint64                 `json:"id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	Config             map[string]interface{} `json:"config"`
	ConfigRaw          []byte                 `json:"config_raw"`
	*AuditableResponse `json:"auditable_response"`
}
