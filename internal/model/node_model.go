package model

type NodeResponse struct {
	Id                uint64                 `json:"id"`
	Type              string                 `json:"type"`
	Label             string                 `json:"label"`
	Description       string                 `json:"description"`
	HelpText          string                 `json:"help_text"`
	Color             string                 `json:"color"`
	Icon              string                 `json:"icon"`
	ComponentName     string                 `json:"component_name"`
	DefaultConfig     map[string]interface{} `json:"default_config"`
	AuditableResponse *AuditableResponse     `json:"auditable_response" mapstructure:"-"`
}

type CreateNodeDto struct {
	Type          string                 `json:"type"`
	Label         string                 `json:"label"`
	Description   string                 `json:"description"`
	HelpText      string                 `json:"help_text"`
	Color         string                 `json:"color"`
	Icon          string                 `json:"icon"`
	ComponentName string                 `json:"component_name"`
	DefaultConfig map[string]interface{} `json:"default_config"`
}

type UpdateNodeDto struct {
	Id            uint64                 `json:"id" validate:"required,number"`
	Type          string                 `json:"type"`
	Label         string                 `json:"label"`
	Description   string                 `json:"description"`
	HelpText      string                 `json:"help_text"`
	Color         string                 `json:"color"`
	Icon          string                 `json:"icon"`
	ComponentName string                 `json:"component_name"`
	DefaultConfig map[string]interface{} `json:"default_config"`
}

type DeleteNodeDto struct {
	Id uint64 `json:"id" validate:"required,number"`
}

type NodeData struct {
	Label string `json:"label"`
	Color string `json:"color"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
}

type Node struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
	Data     NodeData `json:"data"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type Pipeline struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
