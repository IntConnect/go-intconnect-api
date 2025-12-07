package entity

type ModbusServer struct {
	Id        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	IpAddress string    `gorm:"column:ip_address"`
	Port      string    `gorm:"column:port"`
	SlaveId   string    `gorm:"column:slave_id"`
	Timeout   int       `gorm:"column:timeout"`
	IsActive  bool      `gorm:"column:is_active"`
	Auditable Auditable `gorm:"embedded"`
}

func (mqttTopicEntity ModbusServer) GetAuditable() *Auditable {
	return &mqttTopicEntity.Auditable
}
