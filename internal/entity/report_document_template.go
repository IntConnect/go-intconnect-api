package entity

type ReportDocumentTemplate struct {
	Id        uint64       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string       `gorm:"column:name"`
	Code      string       `gorm:"column:code"`
	Parameter []*Parameter `gorm:"many2many:report_document_templates_parameters;joinForeignKey:ReportDocumentTemplateID;joinReferences:ParameterID"`

	Auditable `gorm:"embedded"`
}

func (reportDocumentTemplateEntity ReportDocumentTemplate) GetAuditable() *Auditable {
	return &reportDocumentTemplateEntity.Auditable
}
