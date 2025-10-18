package model

type CreatePipelineEdgeRequest struct {
	SourceNodeTempId string         `json:"source_node_temp_id"`
	TargetNodeTempId string         `json:"target_node_temp_id"`
	Data             map[string]any `json:"data"`
}

type PipelineEdgeResponse struct {
	Id           uint           `json:"id"`
	PipelineId   uint           `json:"pipeline_id"`
	SourceNodeId uint64         `json:"source_node_id"`
	TargetNodeId uint64         `json:"target_node_id"`
	Data         map[string]any `json:"data"`
}
