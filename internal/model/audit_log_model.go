package model

type CreateAuditLogRequest struct {
	UserId      uint64                 `json:"user_id"`
	Action      string                 `json:"action"`
	Feature     string                 `json:"feature"`
	Description string                 `json:"description"`
	Before      map[string]interface{} `json:"before"`
	After       map[string]interface{} `json:"after"`
	IpAddress   string                 `json:"ip_address"`
}

type AuditLogResponse struct {
	Id                uint64                 `json:"id"`
	UserId            uint64                 `json:"user_id"`
	Action            string                 `json:"action"`
	Feature           string                 `json:"feature"`
	Description       string                 `json:"description"`
	Before            map[string]interface{} `json:"before"`
	After             map[string]interface{} `json:"after"`
	IpAddress         string                 `json:"ip_address"`
	AuditableResponse *AuditableResponse     `json:"auditable"`
}

func (auditLogResponse *AuditLogResponse) GetAuditableResponse() *AuditableResponse {
	return auditLogResponse.AuditableResponse
}

func (auditLogResponse *AuditLogResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	auditLogResponse.AuditableResponse = auditableResponse
}
