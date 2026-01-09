package entity

type CheckSheetValue struct {
	Id           uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	CheckSheetId uint64     `gorm:"column:check_sheet_id;"`
	ParameterId  uint64     `gorm:"column:parameter_id;"`
	Timestamp    string     `gorm:"column:timestamp"`
	Value        string     `gorm:"column:value"`
	Parameter    Parameter  `gorm:"foreignKey:ParameterId;references:Id"`
	CheckSheet   CheckSheet `gorm:"foreignKey:CheckSheetId;references:Id"`
}

func (checkSheetValue *CheckSheetValue) GetId() uint64 {
	return checkSheetValue.Id
}
