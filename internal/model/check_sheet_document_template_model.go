package model

type CheckSheetDocumentTemplateResponse struct {
	Id                uint64               `json:"id"`
	Name              string               `json:"name"`
	No                string               `json:"no"`
	Description       string               `json:"description"`
	Category          string               `json:"category"`
	Rotation          int                  `json:"rotation"`
	RotationType      string               `json:"rotation_type"`
	Interval          int                  `json:"interval"`
	RevisionNumber    int                  `json:"revision_number"`
	EffectiveDate     string               `json:"effective_date"`
	ParameterResponse []*ParameterResponse `json:"parameters" mapstructure:"parameters"`
	AuditableResponse *AuditableResponse   `json:"auditable"`
}

type CreateCheckSheetDocumentTemplateRequest struct {
	Name           string   `json:"name" validate:"required,min=3,max=255"`
	No             string   `json:"no" validate:"required,min=3,max=255"`
	Description    string   `json:"description"`
	ParameterIds   []uint64 `json:"parameter_ids" validate:"required,min=1,dive,number,gt=0"`
	Category       string   `json:"category" validate:"required"`
	Rotation       int      `json:"rotation" validate:"required,gte=1"`
	RotationType   string   `json:"rotation_type" validate:"required,oneof=Day Week Month"`
	Interval       int      `json:"interval" validate:"required,gte=1"`
	RevisionNumber int      `json:"revision_number"`
	EffectiveDate  string   `json:"effective_date" validate:"required,date"`
}

type UpdateCheckSheetDocumentTemplateRequest struct {
	Id             uint64   `json:"-" validate:"required,number,gt=0,exists=check_sheet_document_templates;id"`
	Name           string   `json:"name" validate:"required,min=3,max=255"`
	No             string   `json:"no" validate:"required,min=3,max=255"`
	Description    string   `json:"description"`
	ParameterIds   []uint64 `json:"parameter_ids" validate:"required,min=1,dive,number,gt=0"`
	Category       string   `json:"category" validate:"required"`
	Rotation       int      `json:"rotation" validate:"required,gte=1"`
	RotationType   string   `json:"rotation_type" validate:"required,oneof=Day Week Month"`
	Interval       int      `json:"interval" validate:"required,gte=1"`
	RevisionNumber int      `json:"revision_number"`
	EffectiveDate  string   `json:"effective_date" validate:"required,date"`
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetDocumentTemplateResponse.AuditableResponse
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetDocumentTemplateResponse.AuditableResponse = auditableResponse
}
