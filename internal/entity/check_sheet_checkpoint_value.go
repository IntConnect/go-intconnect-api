package entity

type CheckSheetCheckPointValue struct {
	Id                     uint64                `gorm:"column:id;primaryKey;autoIncrement"`
	CheckSheetCheckPointId uint64                `gorm:"column:check_sheet_check_point_id"`
	Timestamp              string                `gorm:"column:timestamp"`
	Value                  string                `gorm:"column:value"`
	CheckSheetCheckPoint   *CheckSheetCheckPoint `gorm:"foreignKey:CheckSheetCheckPointId;references:Id"`
}
