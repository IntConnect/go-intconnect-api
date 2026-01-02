package model

type MqttTopicResponse struct {
	Id                 uint64               `json:"id"`
	MachineId          uint64               `json:"machine_id"`
	MqttBrokerId       uint64               `json:"mqtt_broker_id"`
	Name               string               `json:"name"`
	QoS                int                  `json:"qos"`
	MqttBrokerResponse *MqttBrokerResponse  `json:"mqtt_broker" mapstructure:"MqttBroker"`
	MachineResponse    *MachineResponse     `json:"machine" mapstructure:"Machine"`
	ParameterResponses []*ParameterResponse `json:"parameters" mapstructure:"Parameters"`
	AuditableResponse  *AuditableResponse   `json:"auditable"`
}

type MqttTopicDependency struct {
	MachineResponses    []MachineResponse    `json:"machines"`
	MqttBrokerResponses []MqttBrokerResponse `json:"mqtt_brokers"`
}

type CreateMqttTopicRequest struct {
	MachineId    uint64 `json:"machine_id" validate:"required,number,gte=1,exists=machines;id"`
	MqttBrokerId uint64 `json:"mqtt_broker_id" validate:"required,number,gte=1,exists=mqtt_brokers;id"`
	Name         string `json:"name" validate:"required,min=3,max=255"`
	QoS          int    `json:"qos" validate:"oneof=0 1 2"`
}

type UpdateMqttTopicRequest struct {
	Id           uint64 `json:"-" validate:"required,exists=mqtt_topics;id"`
	MachineId    uint64 `json:"machine_id" validate:"required,number,gte=1,exists=machines;id"`
	MqttBrokerId uint64 `json:"mqtt_broker_id" validate:"required,number,gte=1,exists=mqtt_brokers;id"`
	Name         string `json:"name" validate:"required,min=3,max=255"`
	QoS          int    `json:"qos" validate:"oneof=0 1 2"`
}

type MqttTopicListenerResponse struct {
	SubscribeMultiple SubscribeMultiple
	TopicParameter    TopicParameter
}

type TopicParameter map[string]map[string]uint64
type SubscribeMultiple map[string]byte

func (mqttTopicResponse *MqttTopicResponse) GetAuditableResponse() *AuditableResponse {
	return mqttTopicResponse.AuditableResponse
}

func (mqttTopicResponse *MqttTopicResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	mqttTopicResponse.AuditableResponse = auditableResponse
}
