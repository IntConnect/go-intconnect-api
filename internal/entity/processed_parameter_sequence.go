package entity

type ProcessedParameterSequence struct {
	Id                uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	ParentParameterId uint64     `gorm:"column:parent_parameter_id;"`
	ParameterId       uint64     `gorm:"column:parameter_id;"`
	Sequence          int        `gorm:"column:sequence;"`
	Type              string     `gorm:"column:type;"`
	ParentParameter   *Parameter `gorm:"foreignKey:ParentParameterId;references:Id"`
	Parameter         *Parameter `gorm:"foreignKey:ParameterId;references:Id"`
}
