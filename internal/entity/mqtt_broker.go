package entity

type MqttBroker struct {
	Id       uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	HostName string `gorm:"column:host_name"`
	MqttPort string `gorm:"column:mqtt_port"`
	WsPort   string `gorm:"column:ws_port"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	IsActive bool   `gorm:"column:is_active"`
	Auditable
}

func (mqttBrokerEntity MqttBroker) GetAuditable() *Auditable {
	return &mqttBrokerEntity.Auditable
}
