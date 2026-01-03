package entity

type Register struct {
	Id             uint64        `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId      uint64        `gorm:"column:machine_id"`
	ModbusServerId uint64        `gorm:"column:modbus_server_id"`
	MemoryLocation string        `gorm:"column:memory_location"`
	Name           string        `gorm:"column:name"`
	Description    string        `gorm:"column:description"`
	DataType       string        `gorm:"column:data_type"`
	Machine        *Machine      `gorm:"foreignKey:MachineId;references:Id"`
	ModbusServer   *ModbusServer `gorm:"foreignKey:ModbusServerId;references:Id"`
	Auditable
}

func (registerEntity *Register) GetAuditable() *Auditable {
	return &registerEntity.Auditable
}
