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
	Appearance  map[string]any `json:"appearance"`
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
	Appearance  map[string]any `json:"appearance"`
	Description string         `json:"description"`
}

type PipelineNodeResponse struct {
	Id           uint64                  `json:"id"`
	PipelineId   uint64                  `json:"pipeline_id"`
	NodeId       uint64                  `json:"node_id"`
	Type         string                  `json:"type"`
	Label        string                  `json:"label"`
	PositionX    float64                 `json:"position_x"`
	PositionY    float64                 `json:"position_y"`
	Config       *PipelineNodeConfig     `json:"config" mapstructure:"-"`
	Appearance   *PipelineNodeAppearance `json:"appearance" mapstructure:"-"`
	NodeResponse *NodeResponse           `json:"node_response"`
}
type PipelineNodeConfig struct {
	QoS                           uint64 `json:"qos"`
	Name                          string `json:"name"`
	Topic                         string `json:"topic"`
	Action                        string `json:"action"`
	Output                        string `json:"output"`
	NodeTempId                    string `json:"node_id" mapstructure:"node_id"`
	ProtocolConfigurationId       uint64 `json:"protocol_configuration_id" mapstructure:"protocol_configuration_id"`
	ProtocolConfigurationResponse `json:"protocol_configuration_response" mapstructure:"protocol_configuration_response"`
}

type PipelineNodeAppearance struct {
	BackgroundColor string `json:"background_color" mapstructure:"background_color"`
	Icon            string `json:"icon"`
}
