package entity

import (
	"go-intconnect-api/internal/trait"
	"time"
)

type CheckSheetDocumentTemplate struct {
	Id             uint64                                   `gorm:"column:id;primaryKey;autoIncrement"`
	Name           string                                   `gorm:"column:name"`
	No             string                                   `gorm:"column:no"`
	Description    string                                   `gorm:"column:description"`
	Category       trait.CheckSheetDocumentTemplateCategory `gorm:"column:category"`
	Rotation       int                                      `gorm:"column:rotation"`
	RotationType   string                                   `gorm:"column:rotation_type"`
	Interval       int                                      `gorm:"column:interval"`
	RevisionNumber int                                      `gorm:"column:revision_number"`
	EffectiveDate  time.Time                                `gorm:"column:effective_date"`
	Parameters     []*Parameter                             `gorm:"many2many:check_sheet_document_templates_parameters;joinForeignKey:CheckSheetDocumentTemplateID;joinReferences:ParameterID"`
	Auditable      `gorm:"embedded"`
}

func (checkSheetDocumentTemplateEntity *CheckSheetDocumentTemplate) GetAuditable() *Auditable {
	return &checkSheetDocumentTemplateEntity.Auditable
}
