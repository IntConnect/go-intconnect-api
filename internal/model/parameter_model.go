package model

type ParameterResponse struct {
	Id                uint64             `json:"id"`
	MachineId         *uint64            `json:"machine_id"`
	Name              string             `json:"name"`
	Code              string             `json:"code"`
	Unit              string             `json:"unit"`
	MinValue          float32            `json:"min_value"`
	MaxValue          float32            `json:"max_value"`
	Description       string             `json:"description"`
	MachineResponse   MachineResponse    `json:"machine_response" mapstructure:"-"`
	AuditableResponse *AuditableResponse `json:"auditable_response" mapstructure:"-"`
}

type CreateParameterRequest struct {
	MachineId   *uint64 `json:"machine_id,omitempty" validate:"number"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Code        string  `json:"code" validate:"required,min=3,max=100"`
	Unit        string  `json:"unit" validate:"required,min=3,max=100"`
	MinValue    float32 `json:"min_value,omitempty" validate:""`
	MaxValue    float32 `json:"max_value,omitempty" validate:""`
	Description string  `json:"description,omitempty"`
}

type UpdateParameterRequest struct {
	Id          uint64  `json:"id" validate:"required,number"`
	MachineId   *uint64 `json:"machine_id,omitempty" validate:"number"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Code        string  `json:"code" validate:"required,min=3,max=100"`
	Unit        string  `json:"unit" validate:"required,min=3,max=100"`
	MinValue    float32 `json:"min_value,omitempty" validate:""`
	MaxValue    float32 `json:"max_value,omitempty" validate:""`
	Description string  `json:"description,omitempty"`
}

func (parameterResponse *ParameterResponse) GetAuditableResponse() *AuditableResponse {
	return parameterResponse.AuditableResponse
}

func (parameterResponse *ParameterResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	parameterResponse.AuditableResponse = auditableResponse
}
