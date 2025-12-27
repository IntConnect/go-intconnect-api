package model

type CheckSheetResponse struct {
	Id                           uint64             `json:"id"`
	CheckSheetDocumentTemplateId uint64             `json:"check_sheet_document_template_id"`
	ReportedBy                   uint64             `json:"reported_by"`
	VerifiedBy                   uint64             `json:"verified_by"`
	Date                         string             `json:"date"`
	Note                         string             `json:"note"`
	AuditableResponse            *AuditableResponse `json:"auditable"`
}

type CreateCheckSheetRequest struct {
	CheckSheetDocumentTemplateId uint64             `json:"check_sheet_document_template_id" validate:"required,gte=1,exists=check_sheets;id"`
	CheckSheetValues             []*CheckSheetValue `json:"check_sheet_values" validate:"required,min=1,dive,required"`
}

type CheckSheetValue struct {
	CheckSheetReportDocumentTemplateParameterId uint64 `json:"check_sheet_report_document_template_parameter_id"`
	Value                                       string `json:"value"`
}

type UpdateCheckSheetRequest struct {
	Id                           uint64             `json:"-" validate:"required,gte=1,exists=check_sheets;id"`
	CheckSheetDocumentTemplateId uint64             `json:"check_sheet_document_template_id" validate:"required,gte=1,exists=check_sheet_document_templates;id"`
	CheckSheetValues             []*CheckSheetValue `json:"check_sheet_values" validate:"required,min=1,dive,required"`
}

func (checkSheetResponse *CheckSheetResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetResponse.AuditableResponse
}

func (checkSheetResponse *CheckSheetResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetResponse.AuditableResponse = auditableResponse
}

func (checkSheetValue *CheckSheetValue) GetId() uint64 {
	return checkSheetValue.CheckSheetReportDocumentTemplateParameterId
}
