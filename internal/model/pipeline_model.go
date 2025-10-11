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
