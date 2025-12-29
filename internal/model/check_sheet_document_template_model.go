package model

type CheckSheetDocumentTemplateResponse struct {
	Id                                   uint64                                         `json:"id"`
	Name                                 string                                         `json:"name"`
	No                                   string                                         `json:"no"`
	Description                          string                                         `json:"description"`
	Category                             string                                         `json:"category"`
	Interval                             int                                            `json:"interval"`
	IntervalType                         string                                         `json:"interval_type"`
	RevisionNumber                       int                                            `json:"revision_number"`
	EffectiveDate                        string                                         `json:"effective_date"`
	CheckSheetDocumentTemplateParameters []*CheckSheetDocumentTemplateParameterResponse `json:"parameters" mapstructure:"CheckSheetDocumentTemplateParameters"`
	AuditableResponse                    *AuditableResponse                             `json:"auditable"`
}

type CreateCheckSheetDocumentTemplateRequest struct {
	Name           string   `json:"name" validate:"required,min=3,max=255"`
	No             string   `json:"no" validate:"required,min=3,max=255"`
	Description    string   `json:"description"`
	ParameterIds   []uint64 `json:"parameter_ids" validate:"required,min=1,dive,number,gt=0"`
	Category       string   `json:"category" validate:"required"`
	Interval       int      `json:"interval" validate:"required,gte=1"`
	IntervalType   string   `json:"interval_type" validate:"required,oneof=Hours Minutes"`
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
	Interval       int      `json:"interval" validate:"required,gte=1"`
	IntervalType   string   `json:"interval_type" validate:"required,oneof= Hours Minutes"`
	RevisionNumber int      `json:"revision_number"`
	EffectiveDate  string   `json:"effective_date" validate:"required,date"`
}

type CheckSheetDocumentTemplateParameterResponse struct {
	Id                           uint64             `json:"id"`
	CheckSheetDocumentTemplateId uint64             `json:"check_sheet_document_template_id"`
	ParameterId                  uint64             `json:"parameter_id"`
	ParameterResponse            *ParameterResponse `json:"parameter" mapstructure:"parameter"`
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetDocumentTemplateResponse.AuditableResponse
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetDocumentTemplateResponse.AuditableResponse = auditableResponse
}
