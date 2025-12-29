package model

type CheckSheetResponse struct {
	Id                           uint64                             `json:"id"`
	CheckSheetDocumentTemplateId uint64                             `json:"check_sheet_document_template_id"`
	ReportedBy                   uint64                             `json:"reported_by"`
	VerifiedBy                   uint64                             `json:"verified_by"`
	Timestamp                    string                             `json:"timestamp"`
	Note                         string                             `json:"note"`
	Status                       string                             `json:"status"`
	VerifiedByUser               *UserResponse                      `json:"verified_by_user"`
	ReportedByUser               UserResponse                       `json:"reported_by_user"`
	CheckSheetDocumentTemplate   CheckSheetDocumentTemplateResponse `json:"check_sheet_document_template"`
	CheckSheetValues             []*CheckSheetValueResponse         `json:"check_sheet_values"`
	AuditableResponse            *AuditableResponse                 `json:"auditable"`
}

type CreateCheckSheetRequest struct {
	CheckSheetDocumentTemplateId uint64             `json:"check_sheet_document_template_id" validate:"required,gte=1,exists=check_sheet_document_templates;id"`
	CheckSheetValues             []*CheckSheetValue `json:"check_sheet_values" validate:"required,min=1,dive,required"`
}

type CheckSheetValue struct {
	CheckSheetDocumentTemplateParameterId uint64 `json:"check_sheet_document_template_parameter_id"  validate:"required"`
	Timestamp                             string `json:"timestamp" validate:"required"`
	Value                                 string `json:"value"`
}

type ApprovalCheckSheet struct {
	CheckSheetId uint64 `json:"-" validate:"required,gte=1,exists=check_sheets;id"`
	Note         string `json:"note"`
	Decision     bool   `json:"decision"`
}

type CheckSheetValueResponse struct {
	Id                                    uint64 `json:"id"`
	CheckSheetId                          uint64 `json:"check_sheet_id"`
	CheckSheetDocumentTemplateParameterId uint64 `json:"check_sheet_document_template_parameter_id"`
	Timestamp                             string `json:"timestamp"`
	Value                                 string `json:"value"`
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
	return checkSheetValue.CheckSheetDocumentTemplateParameterId
}
