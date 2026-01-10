package entity

type MqttTopic struct {
	Id           uint64       `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId    uint64       `gorm:"column:machine_id"`
	MqttBrokerId uint64       `gorm:"column:mqtt_broker_id"`
	Name         string       `gorm:"column:name"`
	QoS          int          `gorm:"column:qos"`
	MqttBroker   MqttBroker   `gorm:"foreignKey:MqttBrokerId"`
	Machine      Machine      `gorm:"foreignKey:MachineId;references:Id"`
	Parameters   []*Parameter `gorm:"foreignKey:MqttTopicId;references:Id"`
	Auditable    Auditable    `gorm:"embedded"`
}

func (mqttTopicEntity *MqttTopic) GetAuditable() *Auditable {
	return &mqttTopicEntity.Auditable
}
