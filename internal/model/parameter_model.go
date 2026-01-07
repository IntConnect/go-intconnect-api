package model

type ParameterResponse struct {
	Id                         uint64                               `json:"id"`
	MqttTopicId                *uint64                              `json:"mqtt_topic_id"`
	MachineId                  *uint64                              `json:"machine_id"`
	Name                       string                               `json:"name"`
	Code                       string                               `json:"code"`
	Unit                       string                               `json:"unit"`
	MinValue                   float32                              `json:"min_value"`
	MaxValue                   float32                              `json:"max_value"`
	Description                string                               `json:"description"`
	PositionX                  float32                              `json:"position_x"`
	PositionY                  float32                              `json:"position_y"`
	PositionZ                  float32                              `json:"position_z"`
	RotationX                  float32                              `json:"rotation_x"`
	RotationY                  float32                              `json:"rotation_y"`
	RotationZ                  float32                              `json:"rotation_z"`
	IsAutomatic                bool                                 `json:"is_automatic"`
	IsDisplay                  bool                                 `json:"is_display"`
	IsWatch                    bool                                 `json:"is_watch"`
	IsFeatured                 bool                                 `json:"is_featured"`
	IsRunningTime              bool                                 `json:"is_running_time"`
	IsProcessed                bool                                 `json:"is_processed"`
	MqttTopic                  *MqttTopicResponse                   `json:"mqtt_topic" mapstructure:"MqttTopic"`
	Machine                    *MachineResponse                     `json:"machine" mapstructure:"Machine"`
	ProcessedParameterSequence []ProcessedParameterSequenceResponse `json:"processed_parameter_sequence" mapstructure:"-"`
	ParameterOperation         []ParameterOperationRequest          `json:"parameter_operation" mapstructure:"-c"`
	AuditableResponse          *AuditableResponse                   `json:"auditable"`
}

type ParameterDependency struct {
	MachineResponses   []MachineResponse   `json:"machines"`
	MqttTopicResponses []MqttTopicResponse `json:"mqtt_topics"`
}

type CreateParameterRequest struct {
	MachineId                  *uint64                              `json:"machine_id" validate:"omitempty,number,gte=1,exists=machines;id" property:"Machine"`
	MqttTopicId                *uint64                              `json:"mqtt_topic_id" validate:"omitempty,number,gte=1,exists=mqtt_topics;id" property:"MQTT Topic"`
	Name                       string                               `json:"name" validate:"required,min=3,max=100"`
	Code                       string                               `json:"code" validate:"required,min=3,max=100"`
	Unit                       string                               `json:"unit" validate:"required,min=1,max=100"`
	MinValue                   *float32                             `json:"min_value,omitempty"`
	MaxValue                   *float32                             `json:"max_value,omitempty"`
	Description                string                               `json:"description"`
	Category                   string                               `json:"category"`
	PositionX                  *float32                             `json:"position_x,omitempty"`
	PositionY                  *float32                             `json:"position_y,omitempty"`
	PositionZ                  *float32                             `json:"position_z,omitempty"`
	RotationX                  *float32                             `json:"rotation_x,omitempty"`
	RotationY                  *float32                             `json:"rotation_y,omitempty"`
	RotationZ                  *float32                             `json:"rotation_z,omitempty"`
	IsDisplay                  bool                                 `json:"is_display"`
	IsAutomatic                bool                                 `json:"is_automatic"`
	IsWatch                    bool                                 `json:"is_watch"`
	IsFeatured                 bool                                 `json:"is_featured"`
	IsRunningTime              bool                                 `json:"is_running_time"`
	IsProcessed                bool                                 `json:"is_processed"`
	ProcessedParameterSequence []*ProcessedParameterSequenceRequest `json:"processed_parameter_sequence" validate:"omitempty"`
}

type UpdateParameterRequest struct {
	Id                         uint64                               `json:"-" validate:"required,number,gte=1,exists=parameters;id"`
	MachineId                  *uint64                              `json:"machine_id" validate:"omitempty,number,gte=1,exists=machines;id" property:"Machine"`
	MqttTopicId                *uint64                              `json:"mqtt_topic_id" validate:"omitempty,required,number,gte=1,exists=mqtt_topics;id" property:"MQTT Topic"`
	Name                       string                               `json:"name" validate:"required,min=3,max=100"`
	Code                       string                               `json:"code" validate:"required,min=3,max=100"`
	Unit                       string                               `json:"unit" validate:"required,min=1,max=100"`
	MinValue                   *float32                             `json:"min_value,omitempty"`
	MaxValue                   *float32                             `json:"max_value,omitempty"`
	Description                string                               `json:"description"`
	PositionX                  *float32                             `json:"position_x,omitempty"`
	PositionY                  *float32                             `json:"position_y,omitempty"`
	PositionZ                  *float32                             `json:"position_z,omitempty"`
	RotationX                  *float32                             `json:"rotation_x,omitempty"`
	RotationY                  *float32                             `json:"rotation_y,omitempty"`
	RotationZ                  *float32                             `json:"rotation_z,omitempty"`
	IsDisplay                  bool                                 `json:"is_display"`
	IsAutomatic                bool                                 `json:"is_automatic"`
	IsWatch                    bool                                 `json:"is_watch"`
	IsFeatured                 bool                                 `json:"is_featured"`
	IsRunningTime              bool                                 `json:"is_running_time"`
	IsProcessed                bool                                 `json:"is_processed"`
	ProcessedParameterSequence []*ProcessedParameterSequenceRequest `json:"processed_parameter_sequence" validate:"omitempty"`
}

type ManageParameterOperationRequest struct {
	Id      uint64                       `json:"-" validate:"required,number"`
	Created []*ParameterOperationRequest `json:"created"`
	Updated []*ParameterOperationRequest `json:"updated"`
	Deleted []uint64                     `json:"deleted"`
}

type ParameterOperationRequest struct {
	Type     string  `json:"type" validate:"required,oneof=ADDITION SUBTRACTION MULTIPLICATION DIVISION"`
	Value    float32 `json:"value" validate:"required"`
	Sequence int     `json:"sequence" validate:"required"`
}
type ParameterOperationResponse struct {
	Id       uint64  `json:"id"`
	Type     string  `json:"type"`
	Value    float32 `json:"value"`
	Sequence int     `json:"sequence"`
}

type ParameterFilterRequest struct {
	IsAutomatic *string `form:"omitempty,is_automatic"`
	IsDisplay   *string `form:"omitempty,is_display"`
	IsWatch     *string `form:"omitempty,is_watch"`
	IsFeatured  *string `form:"omitempty,is_featured"`
	IsProcessed *string `form:"omitempty,is_processed"`
}

type ProcessedParameterSequenceRequest struct {
	ParameterId uint64 `json:"parameter_id" validate:"required,gte=1,exists=parameters;id"`
	Sequence    int    `json:"sequence" validate:"required;gte=1"`
	Type        string `json:"type" validate:"required,oneof=ADDITION MULTIPLICATION DIVISION"`
}

type ProcessedParameterSequenceResponse struct {
	Id                uint64 `json:"id"`
	ParentParameterId uint64 `json:"parent_parameter_id"`
	ParameterId       uint64 `json:"parameter_id"`
	Sequence          int    `json:"sequence"`
	Type              string `json:"type"`
}

func (parameterResponse *ParameterResponse) GetAuditableResponse() *AuditableResponse {
	return parameterResponse.AuditableResponse
}

func (parameterResponse *ParameterResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	parameterResponse.AuditableResponse = auditableResponse
}
