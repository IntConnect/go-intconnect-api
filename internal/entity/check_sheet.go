package entity

import "time"

type CheckSheet struct {
	Id                           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	CheckSheetDocumentTemplateId uint64    `gorm:"column:check_sheet_document_template_id"`
	ReportedBy                   uint64    `gorm:"column:reported_by"`
	VerifiedBy                   uint64    `gorm:"column:verified_by"`
	Date                         time.Time `gorm:"column:date"`
	Note                         string    `gorm:"column:note"`
	Auditable                    `gorm:"embedded"`
}
