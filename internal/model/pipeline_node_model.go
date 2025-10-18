package model

type CreatePipelineNodeRequest struct {
	TempID      string         `json:"temp_id"` // id VueFlow
	NodeID      uint64         `json:"node_id"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Label       string         `json:"label"`
	PositionX   float64        `json:"position_x"`
	PositionY   float64        `json:"position_y"`
	Config      map[string]any `json:"config"`
	Description string         `json:"description"`
}

type UpdatePipelineNodeRequest struct {
	Id          uint64         `json:"id"`
	TempID      string         `json:"temp_id"` // id VueFlow
	NodeID      uint64         `json:"node_id"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Label       string         `json:"label"`
	PositionX   float64        `json:"position_x"`
	PositionY   float64        `json:"position_y"`
	Config      map[string]any `json:"config"`
	Description string         `json:"description"`
}

type PipelineNodeResponse struct {
	Id           uint64         `json:"id"`
	PipelineId   uint64         `json:"pipeline_id"`
	NodeId       uint64         `json:"node_id"`
	Type         string         `json:"type"`
	Label        string         `json:"label"`
	PositionX    float64        `json:"position_x"`
	PositionY    float64        `json:"position_y"`
	Config       map[string]any `json:"config"`
	NodeResponse *NodeResponse  `json:"node_response"`
}
