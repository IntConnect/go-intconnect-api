package entity

type CheckSheetDocumentTemplateParameter struct {
	Id                           uint64                      `gorm:"column:id;primaryKey;autoIncrement"`
	CheckSheetDocumentTemplateId uint64                      `gorm:"column:check_sheet_document_template_id;"`
	ParameterId                  uint64                      `gorm:"column:parameter_id;"`
	CheckSheetDocumentTemplate   *CheckSheetDocumentTemplate `gorm:"foreignKey:CheckSheetDocumentTemplateId;references:Id"`
	Parameter                    *Parameter                  `gorm:"foreignKey:ParameterId;references:Id"`
}

func (CheckSheetDocumentTemplateParameter) TableName() string {
	return "check_sheet_document_templates_parameters"
}

func (checkSheetDocumentTemplateParameter *CheckSheetDocumentTemplateParameter) GetId() uint64 {
	return checkSheetDocumentTemplateParameter.Id
}
