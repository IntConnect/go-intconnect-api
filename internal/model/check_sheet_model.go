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
	CheckSheetCheckPoints        []*CheckSheetCheckPointResponse    `json:"check_sheet_check_points"`
	AuditableResponse            *AuditableResponse                 `json:"auditable"`
}

type CreateCheckSheetRequest struct {
	CheckSheetDocumentTemplateId uint64                  `json:"check_sheet_document_template_id" validate:"required,gte=1,exists=check_sheet_document_templates;id"`
	CheckSheetCheckPoints        []*CheckSheetCheckPoint `json:"check_sheet_check_points" validate:"required,dive"`
}

type CheckSheetCheckPoint struct {
	Id               uint64             `json:"id"`
	CheckSheetId     uint64             `json:"check_sheet_id"`
	ParameterId      uint64             `json:"parameter_id"`
	Name             string             `json:"name"`
	CheckSheetValues []*CheckSheetValue `json:"check_sheet_values" validate:"required,min=1,dive,required"`
}

type CheckSheetCheckPointResponse struct {
	Id                        uint64                               `json:"id"`
	CheckSheetId              uint64                               `json:"check_sheet_id;"`
	ParameterId               uint64                               `json:"parameter_id"`
	Name                      string                               `json:"name"`
	CheckSheetCheckPointValue []*CheckSheetCheckPointValueResponse `json:"check_sheet_check_point_value"`
}

type CheckSheetValue struct {
	Timestamp string `json:"timestamp" validate:"required"`
	Value     string `json:"value"`
}

type CheckSheetCheckPointValueResponse struct {
	Id        uint64 `json:"id"`
	Timestamp string `json:"timestamp" `
	Value     string `json:"value"`
}

type ApprovalCheckSheet struct {
	CheckSheetId uint64 `json:"-" validate:"required,gte=1,exists=check_sheets;id"`
	Note         string `json:"note"`
	Decision     bool   `json:"decision"`
}

type CheckSheetValueResponse struct {
	Id           uint64 `json:"id"`
	CheckSheetId uint64 `json:"check_sheet_id"`
	ParameterId  uint64 `json:"parameter_id"`
	Timestamp    string `json:"timestamp"`
	Value        string `json:"value"`
}

type UpdateCheckSheetRequest struct {
	Id                           uint64                  `json:"-" validate:"required,gte=1,exists=check_sheets;id"`
	CheckSheetDocumentTemplateId uint64                  `json:"check_sheet_document_template_id" validate:"required,gte=1,exists=check_sheet_document_templates;id"`
	CheckSheetCheckPoint         []*CheckSheetCheckPoint `json:"check_sheet_check_point" validate:"required,gt=0"`
	CheckSheetValues             []*CheckSheetValue      `json:"check_sheet_values" validate:"required,min=1,dive,required"`
}

func (checkSheetResponse *CheckSheetResponse) GetAuditableResponse() *AuditableResponse {
	return checkSheetResponse.AuditableResponse
}

func (checkSheetResponse *CheckSheetResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	checkSheetResponse.AuditableResponse = auditableResponse
}

func (checkSheetCheckPoint *CheckSheetCheckPoint) GetId() uint64 {
	return checkSheetCheckPoint.Id
}
