package model

type MqttTopicResponse struct {
	Id           uint64 `json:"id"`
	MqttBrokerId uint64 `gorm:"column:mqtt_broker_id"`
	Name         string `gorm:"column:name"`
	QoS          int    `gorm:"column:qos"`
}

type CreateMqttTopicRequest struct {
	MqttBrokerId uint64 `gorm:"column:mqtt_broker_id"`
	Name         string `gorm:"column:name"`
	QoS          int    `gorm:"column:qos"`
}

type UpdateMqttTopicRequest struct {
	Id           uint64 `json:"id"`
	MqttBrokerId uint64 `gorm:"column:mqtt_broker_id"`
	Name         string `gorm:"column:name"`
	QoS          int    `gorm:"column:qos"`
}

type MqttTopicListenerResponse struct {
	SubscribeMultiple SubscribeMultiple
	TopicParameter    TopicParameter
}

type TopicParameter map[string]map[string]float64
type SubscribeMultiple map[string]byte
