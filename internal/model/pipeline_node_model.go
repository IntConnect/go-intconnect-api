package model

type CreatePipelineNodeDto struct {
	TempID      string         `json:"temp_id"` // id VueFlow
	NodeID      uint64         `json:"node_id"`
	Type        string         `json:"type"`
	Label       string         `json:"label"`
	PositionX   float64        `json:"position_x"`
	PositionY   float64        `json:"position_y"`
	Config      map[string]any `json:"config"`
	Description string         `json:"description"`
}

type UpdatePipelineNodeDto struct {
	Id          uint64                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Config      map[string]interface{} `json:"config"`
}
