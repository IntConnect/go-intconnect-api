package model

type CheckSheetDocumentTemplateResponse struct {
	Id                uint64               `json:"id"`
	Name              string               `json:"name" validate:"required,min=3,max=255"`
	Code              string               `json:"code" validate:"required,min=3,max=255"`
	DocumentVersion   int                  `json:"document_version"`
	ParameterResponse []*ParameterResponse `json:"parameters" mapstructure:"parameters"`
	AuditableResponse *AuditableResponse   `json:"auditable"`
}

type CreateCheckSheetDocumentTemplateRequest struct {
	Name         string   `json:"name" validate:"required,min=3,max=255"`
	Code         string   `json:"code" validate:"required,min=3,max=255"`
	ParameterIds []uint64 `json:"parameter_ids" validate:"required,min=1,dive,number,gt=0"`
}

type UpdateCheckSheetDocumentTemplateRequest struct {
	Id           uint64   `json:"-" validate:"required,number,gt=0,exists=check_sheet_document_templates;id"`
	Name         string   `json:"name" validate:"required,min=3,max=255"`
	Code         string   `json:"code" validate:"required,min=3,max=255"`
	ParameterIds []uint64 `json:"parameter_ids" validate:"required,min=1,dive,exists=parameters;id"`
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetDocumentTemplateResponse.AuditableResponse
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetDocumentTemplateResponse.AuditableResponse = auditableResponse
}
