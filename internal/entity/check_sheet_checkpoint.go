package entity

type CheckSheetCheckPoint struct {
	Id                        uint64                       `gorm:"column:id;primaryKey;autoIncrement"`
	CheckSheetId              uint64                       `gorm:"column:check_sheet_id"`
	ParameterId               uint64                       `gorm:"column:parameter_id"`
	Name                      string                       `gorm:"column:name"`
	CheckSheetCheckPointValue []*CheckSheetCheckPointValue `gorm:"foreignKey:CheckSheetCheckPointId;references:Id"`
	CheckSheet                CheckSheet                   `gorm:"foreignKey:CheckSheetId;references:Id"`
	Auditable                 Auditable                    `gorm:"embedded"`
}

func (checkSheetCheckPoint *CheckSheetCheckPoint) GetAuditable() *Auditable {
	return &checkSheetCheckPoint.Auditable
}

func (checkSheetCheckPoint *CheckSheetCheckPoint) GetId() uint64 {
	return checkSheetCheckPoint.Id
}
