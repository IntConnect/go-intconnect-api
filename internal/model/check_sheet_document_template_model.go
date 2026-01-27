package model

type CheckSheetDocumentTemplateResponse struct {
	Id                uint64             `json:"id"`
	MachineId         uint64             `json:"machine_id"`
	Name              string             `json:"name"`
	No                string             `json:"no"`
	Description       string             `json:"description"`
	Category          string             `json:"category"`
	Interval          int                `json:"interval"`
	IntervalType      string             `json:"interval_type"`
	RotationType      string             `json:"rotation_type"`
	RevisionNumber    int                `json:"revision_number"`
	EffectiveDate     string             `json:"effective_date"`
	StartingHour      string             `json:"starting_hour"`
	AuditableResponse *AuditableResponse `json:"auditable"`
}

type CreateCheckSheetDocumentTemplateRequest struct {
	MachineId    uint64 `json:"machine_id" validate:"required,gte=1,exists=machines;id"`
	Name         string `json:"name" validate:"required,min=3,max=255"`
	No           string `json:"no" validate:"required,min=3,max=255"`
	Description  string `json:"description"`
	Category     string `json:"category" validate:"required"`
	Interval     int    `json:"interval" validate:"required,gte=1"`
	RotationType string `json:"rotation_type" validate:"required,oneof=Daily Weekly Monthly"`
	IntervalType string `json:"interval_type" validate:"required,oneof=Hour Minute"`

	RevisionNumber int    `json:"revision_number"`
	EffectiveDate  string `json:"effective_date" validate:"required,date"`
	StartingHour   string `json:"starting_hour" validate:"required,time"`
}

type UpdateCheckSheetDocumentTemplateRequest struct {
	Id             uint64 `json:"-" validate:"required,number,gte=1,exists=check_sheet_document_templates;id"`
	MachineId      uint64 `json:"machine_id" validate:"required,gte=1,exists=machines;id"`
	Name           string `json:"name" validate:"required,min=3,max=255"`
	No             string `json:"no" validate:"required,min=3,max=255"`
	Description    string `json:"description"`
	Category       string `json:"category" validate:"required"`
	Interval       int    `json:"interval" validate:"required,gte=1"`
	IntervalType   string `json:"interval_type" validate:"required,oneof=Hour Minute"`
	RotationType   string `json:"rotation_type" validate:"required,oneof=Daily Weekly Monthly"`
	RevisionNumber int    `json:"revision_number"`
	EffectiveDate  string `json:"effective_date" validate:"required,date"`
	StartingHour   string `json:"starting_hour" validate:"required,time"`
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetDocumentTemplateResponse.AuditableResponse
}

func (checkSheetDocumentTemplateResponse *CheckSheetDocumentTemplateResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetDocumentTemplateResponse.AuditableResponse = auditableResponse
}
