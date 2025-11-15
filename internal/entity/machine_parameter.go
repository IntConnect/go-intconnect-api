package entity

type MachineParameter struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ParameterId uint64    `gorm:"column:parameter_id;"`
	MeshName    string    `gorm:"column:mesh_name"`
	PositionX   float32   `gorm:"column:position_x"`
	PositionY   float32   `gorm:"column:position_y"`
	PositionZ   float32   `gorm:"column:position_z"`
	LabelType   string    `gorm:"column:label_type"`
	Parameter   Parameter `gorm:"foreignKey:ParameterId;references:Id"`
	Auditable
}
