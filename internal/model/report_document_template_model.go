package model

type ReportDocumentTemplateResponse struct {
	Id                uint64              `json:"id"`
	Name              string              `json:"name" validate:"required,min=3,max=255"`
	Code              string              `json:"code" validate:"required,min=3,max=255"`
	ParameterResponse []ParameterResponse `json:"parameters"`
	AuditableResponse *AuditableResponse  `json:"auditable"`
}

type CreateReportDocumentTemplateRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=255"`
	Code        string   `json:"code" validate:"required,min=3,max=255"`
	ParameterId []uint64 `json:"parameter_id" validate:"required,min=1,dive,number,gt=0"`
}

type UpdateReportDocumentTemplateRequest struct {
	Id          uint64   `json:"id" validate:"required,number,gt=0,exists=report_document_templates;id"`
	Name        string   `json:"name" validate:"required,min=3,max=255"`
	Code        string   `json:"code" validate:"required,min=3,max=255"`
	ParameterId []uint64 `json:"parameter_id" validate:"required,min=1,dive,exists=parameters;id"`
}

func (reportDocumentTemplateResponse *ReportDocumentTemplateResponse) GetAuditableResponse() *AuditableResponse {
	return reportDocumentTemplateResponse.AuditableResponse
}

func (reportDocumentTemplateResponse *ReportDocumentTemplateResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	reportDocumentTemplateResponse.AuditableResponse = auditableResponse
}
