package entity

import (
	"go-intconnect-api/internal/trait"
	"time"
)

type CheckSheetDocumentTemplate struct {
	Id             uint64                                   `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId      uint64                                   `gorm:"column:machine_id;"`
	Name           string                                   `gorm:"column:name"`
	No             string                                   `gorm:"column:no"`
	Description    string                                   `gorm:"column:description"`
	Category       trait.CheckSheetDocumentTemplateCategory `gorm:"column:category"`
	Interval       int                                      `gorm:"column:interval"`
	IntervalType   string                                   `gorm:"column:interval_type"`
	RotationType   string                                   `gorm:"column:rotation_type"`
	RevisionNumber int                                      `gorm:"column:revision_number"`
	EffectiveDate  time.Time                                `gorm:"column:effective_date"`
	StartingHour   string                                   `gorm:"column:starting_hour"`
	Machine        Machine                                  `gorm:"foreignKey:MachineId;references:Id"`
	Auditable      Auditable                                `gorm:"embedded"`
}

func (checkSheetDocumentTemplateEntity *CheckSheetDocumentTemplate) GetAuditable() *Auditable {
	return &checkSheetDocumentTemplateEntity.Auditable
}
