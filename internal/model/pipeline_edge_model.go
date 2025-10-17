package model

type CreatePipelineEdgeDto struct {
	SourceNodeTempId string         `json:"source_node_temp_id"`
	TargetNodeTempId string         `json:"target_node_temp_id"`
	Data             map[string]any `json:"data"`
}
