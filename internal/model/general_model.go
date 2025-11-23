package model

type AuditableResponse struct {
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt string `json:"updated_at"`
	DeletedBy string `json:"deleted_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

type HasAuditableResponse interface {
	GetAuditableResponse() *AuditableResponse
	SetAuditableResponse(auditableResponse *AuditableResponse)
}

type DeleteResourceGeneralRequest struct {
	Id uint64 `json:"id" validate:"required|number"`
}

// BaseResponse adalah struktur dasar untuk semua response
type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SuccessResponse untuk response sukses dengan data
type SuccessResponse[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

type ErrorResponse struct {
	BaseResponse
	Error *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail berisi detail error
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
