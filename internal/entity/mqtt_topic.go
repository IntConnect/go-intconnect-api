package entity

type MqttTopic struct {
	Id           uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	MqttBrokerId uint64     `gorm:"column:mqtt_broker_id"`
	Name         string     `gorm:"column:name"`
	QoS          int        `gorm:"column:qos"`
	MqttBroker   MqttBroker `gorm:"foreignKey:MqttBrokerId"`
	Auditable
}

func (mqttTopicEntity MqttTopic) GetAuditable() *Auditable {
	return &mqttTopicEntity.Auditable
}
