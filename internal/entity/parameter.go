package entity

type Parameter struct {
	Id          uint64   `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId   *uint64  `gorm:"column:machine_id;"`
	Name        string   `gorm:"column:name"`
	Code        string   `gorm:"column:code"`
	Unit        string   `gorm:"column:unit"`
	MinValue    float32  `gorm:"column:min_value"`
	MaxValue    float32  `gorm:"column:max_value"`
	Description string   `gorm:"column:description"`
	Machine     *Machine `gorm:"foreignKey:MachineId;references:Id"`
	Auditable
}

func (parameterEntity Parameter) GetAuditable() *Auditable {
	return &parameterEntity.Auditable
}
