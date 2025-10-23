package model

type CreatePipelineRequest struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Config      map[string]interface{}      `json:"config"`
	Nodes       []CreatePipelineNodeRequest `json:"nodes" mapstructure:"-"`
	Edges       []CreatePipelineEdgeRequest `json:"edges" mapstructure:"-"`
}

type UpdatePipelineRequest struct {
	Id          uint64                      `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Config      map[string]interface{}      `json:"config"`
	Nodes       []CreatePipelineNodeRequest `json:"nodes" mapstructure:"-"`
	Edges       []CreatePipelineEdgeRequest `json:"edges" mapstructure:"-"`
}

type PipelineResponse struct {
	Id           uint64                  `json:"id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Config       map[string]interface{}  `json:"config"`
	PipelineNode []*PipelineNodeResponse `json:"pipeline_node"`
	PipelineEdge []*PipelineEdgeResponse `json:"pipeline_edge"`

	*AuditableResponse `json:"auditable_response"`
}

type DeletePipelineRequest struct {
	ID uint64 `json:"id"`
}
