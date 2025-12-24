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
	CheckSheetDocumentTemplateId uint64 `json:"check_sheet_document_template_id" validate:"required,gte=1,exists=check_sheets;id"`
	Date                         string `json:"date" validate:"required,date"`
}

func (checkSheetResponse *CheckSheetResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetResponse.AuditableResponse
}

func (checkSheetResponse *CheckSheetResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetResponse.AuditableResponse = auditableResponse
}
