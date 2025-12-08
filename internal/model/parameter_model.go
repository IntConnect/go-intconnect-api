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
	MachineResponse   *MachineResponse   `json:"machine_response" mapstructure:"-"`
	AuditableResponse *AuditableResponse `json:"auditable_response" mapstructure:"-"`
}

type ParameterDependency struct {
	MachineResponses   []MachineResponse   `json:"machines"`
	MqttTopicResponses []MqttTopicResponse `json:"mqtt_topics"`
}

type CreateParameterRequest struct {
	MachineId   uint64  `json:"machine_id" validate:"required,number,gt=0,exists=machines;id"`
	MqttTopicId uint64  `json:"mqtt_topic_id" validate:"required,number,gt=0,exists=mqtt_topics;id"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Code        string  `json:"code" validate:"required,min=3,max=100"`
	Unit        string  `json:"unit" validate:"required,min=1,max=100"`
	MinValue    float32 `json:"min_value,omitempty"`
	MaxValue    float32 `json:"max_value,omitempty"`
	Description string  `json:"description,omitempty"`
	PositionX   float32 `json:"position_x,omitempty"`
	PositionY   float32 `json:"position_y,omitempty"`
	PositionZ   float32 `json:"position_z,omitempty"`
	RotationX   float32 `json:"rotation_x,omitempty"`
	RotationY   float32 `json:"rotation_y,omitempty"`
	RotationZ   float32 `json:"rotation_z,omitempty"`
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

type ManageParameterOperationRequest struct {
	Id                  uint64                       `json:"id" validate:"required,number"`
	ParameterOperations []*ParameterOperationRequest `json:"parameter_operations" validate:"required,dive,min=1"`
}

type ParameterOperationRequest struct {
	Type     string  `json:"type" validate:"required,oneof=ADDITION SUBTRACTION MULTIPLICATION DIVISION"`
	Value    float32 `json:"value" validate:"required,number"`
	Sequence int     `json:"sequence" validate:"required,number"`
}

func (parameterResponse *ParameterResponse) GetAuditableResponse() *AuditableResponse {
	return parameterResponse.AuditableResponse
}

func (parameterResponse *ParameterResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	parameterResponse.AuditableResponse = auditableResponse
}
