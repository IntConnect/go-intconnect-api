package model

type MqttTopicResponse struct {
	Id                 uint64             `json:"id"`
	Name               string             `json:"name"`
	QoS                int                `json:"qos"`
	MqttBrokerResponse MqttBrokerResponse `json:"mqtt_broker"`
	MachineResponse    MachineResponse    `json:"machine"`
	AuditableResponse  *AuditableResponse `json:"auditable"`
}

type CreateMqttTopicRequest struct {
	MachineId    uint64 `json:"machine_id" validate:"required,number,gt=0,exists=machines;id"`
	MqttBrokerId uint64 `json:"mqtt_broker_id" validate:"required,number,gt=0,exists=mqtt_brokers;id"`
	Name         string `json:"name" validate:"required,min=3,max=255"`
	QoS          int    `json:"qos" validate:"oneof=0 1 2"`
}

type UpdateMqttTopicRequest struct {
	Id           uint64 `json:"id" validate:"required,exists=mqtt_topics;id"`
	MqttBrokerId uint64 `json:"mqtt_broker_id" validate:"required,exists=mqtt_brokers;id"`
	Name         string `json:"name" validate:"required,min=3,max=255"`
	QoS          int    `json:"qos" validate:"required,oneof= 0 1 2"`
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
