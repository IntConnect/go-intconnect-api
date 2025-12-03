package model

type CreateAuditLogRequest struct {
	UserId      uint64 `json:"user_id"`
	Action      string `json:"action"`
	Feature     string `json:"feature"`
	Description string `json:"description"`
	Before      any    `json:"before"`
	After       any    `json:"after"`
	IpAddress   string `json:"ip_address"`
}

type AuditLogResponse struct {
	Id                      uint64                   `json:"id"`
	UserId                  uint64                   `json:"user_id"`
	Action                  string                   `json:"action"`
	Feature                 string                   `json:"feature"`
	Description             string                   `json:"description"`
	Before                  map[string]interface{}   `json:"before"`
	After                   map[string]interface{}   `json:"after"`
	IpAddress               string                   `json:"ip_address"`
	UserResponse            *UserResponse            `json:"user,omitempty" mapstructure:"user"`
	SimpleAuditableResponse *SimpleAuditableResponse `json:"auditable"`
}

func (auditLogResponse *AuditLogResponse) GetSimpleAuditableResponse() *SimpleAuditableResponse {
	return auditLogResponse.SimpleAuditableResponse
}

func (auditLogResponse *AuditLogResponse) SetSimpleAuditableResponse(simpleAuditableResponse *SimpleAuditableResponse) {
	auditLogResponse.SimpleAuditableResponse = simpleAuditableResponse
}

const (
	AUDIT_LOG_LOGIN  = "LOGIN"
	AUDIT_LOG_CREATE = "CREATE"
	AUDIT_LOG_UPDATE = "UPDATE"
	AUDIT_LOG_DELETE = "DELETE"
)

const (
	AUDIT_LOG_FEATURE_USER        = "USER"
	AUDIT_LOG_FEATURE_ROLE        = "ROLE"
	AUDIT_LOG_FEATURE_MQTT_BROKER = "MQTT_BROKER"
)

const (
	AUDIT_LOG_ACTOR_SYSTEM = "SYSTEM"
)
