package model

type PermissionResponse struct {
	Id          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	AuditableResponse
}
