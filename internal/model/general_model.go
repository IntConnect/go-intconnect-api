package model

type AuditableResponse struct {
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt string `json:"updated_at"`
	DeletedBy string `json:"deleted_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

type DeleteResourceGeneralRequest struct {
	Id uint64 `json:"id" validate:"required|number"`
}
