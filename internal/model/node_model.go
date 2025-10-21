package model

type NodeResponse struct {
	Id                uint64                 `json:"id"`
	Type              string                 `json:"type"`
	Name              string                 `json:"name"`
	Label             string                 `json:"label"`
	Description       string                 `json:"description"`
	HelpText          string                 `json:"help_text"`
	Color             string                 `json:"color"`
	Icon              string                 `json:"icon"`
	ComponentName     string                 `json:"component_name"`
	DefaultConfig     map[string]interface{} `json:"default_config" mapstructure:"-"`
	AuditableResponse *AuditableResponse     `json:"auditable_response" mapstructure:"-"`
}

type CreateNodeRequest struct {
	Type          string                 `json:"type"`
	Label         string                 `json:"label"`
	Description   string                 `json:"description"`
	HelpText      string                 `json:"help_text"`
	Color         string                 `json:"color"`
	Icon          string                 `json:"icon"`
	ComponentName string                 `json:"component_name"`
	DefaultConfig map[string]interface{} `json:"default_config"`
}

type UpdateNodeRequest struct {
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

type DeleteNodeRequest struct {
	Id uint64 `json:"id" validate:"required,number"`
}
