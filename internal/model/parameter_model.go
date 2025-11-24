package model

type ParameterResponse struct {
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

type CreateParameterRequest struct {
	MachineId   uint64  `json:"machine_id" validate:"required,number,gt=0"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Code        string  `json:"code" validate:"required,min=3,max=100"`
	Unit        string  `json:"unit" validate:"required,min=3,max=100"`
	MinValue    float32 `json:"min_value"`
	MaxValue    float32 `json:"max_value"`
	Description string  `json:"description"`
}

type UpdateParameterRequest struct {
	Id          uint64  `json:"id" validate:"required,number"`
	MachineId   uint64  `json:"machine_id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Unit        string  `json:"unit"`
	MinValue    float32 `json:"min_value"`
	MaxValue    float32 `json:"max_value"`
	Description string  `json:"description"`
}

type DeleteParameterRequest struct {
	Id uint64 `json:"id" validate:"required,number"`
}

func (parameterResponse *ParameterResponse) GetAuditableResponse() *AuditableResponse {
	return parameterResponse.AuditableResponse
}

func (parameterResponse *ParameterResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	parameterResponse.AuditableResponse = auditableResponse
}
