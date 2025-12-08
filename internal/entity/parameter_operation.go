package entity

type ParameterOperation struct {
	Id          uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	ParameterId uint64     `gorm:"column:parameterId;"`
	Type        string     `gorm:"column:type;"`
	Value       float32    `gorm:"column:value"`
	Sequence    int        `gorm:"column:sequence"`
	Parameter   *Parameter `gorm:"foreignKey:ParameterId;references:Id"`
}
