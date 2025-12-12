package entity

type CheckSheetDocumentTemplate struct {
	Id              uint64       `gorm:"column:id;primaryKey;autoIncrement"`
	Name            string       `gorm:"column:name"`
	Code            string       `gorm:"column:code"`
	DocumentVersion int          `gorm:"column:document_version"`
	Parameters      []*Parameter `gorm:"many2many:check_sheet_document_templates_parameters;joinForeignKey:CheckSheetDocumentTemplateID;joinReferences:ParameterID"`
	Auditable       `gorm:"embedded"`
}

func (checkSheetDocumentTemplateEntity *CheckSheetDocumentTemplate) GetAuditable() *Auditable {
	return &checkSheetDocumentTemplateEntity.Auditable
}
