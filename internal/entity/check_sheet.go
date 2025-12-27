package entity

import (
	"time"
)

type CheckSheet struct {
	Id                           uint64             `gorm:"column:id;primaryKey;autoIncrement"`
	CheckSheetDocumentTemplateId uint64             `gorm:"column:check_sheet_document_template_id"`
	ReportedBy                   uint64             `gorm:"column:reported_by"`
	VerifiedBy                   uint64             `gorm:"column:verified_by"`
	Timestamp                    time.Time          `gorm:"column:timestamp"`
	Note                         string             `gorm:"column:note"`
	CheckSheetValue              []*CheckSheetValue `gorm:"foreignKey:CheckSheetId;references:Id"`
	Auditable                    `gorm:"embedded"`
}

func (checkSheetEntity *CheckSheet) GetAuditable() *Auditable {
	return &checkSheetEntity.Auditable
}
