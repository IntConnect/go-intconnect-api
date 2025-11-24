package model

type AuditableResponse struct {
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt string `json:"updated_at"`
	DeletedBy string `json:"deleted_by"`
	DeletedAt string `json:"deleted_at"`
}

type HasAuditableResponse interface {
	GetAuditableResponse() *AuditableResponse
	SetAuditableResponse(auditableResponse *AuditableResponse)
}
