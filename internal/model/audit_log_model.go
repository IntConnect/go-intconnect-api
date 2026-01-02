package model

type CreateAuditLogRequest struct {
	UserId      uint64 `json:"user_id"`
	Action      string `json:"action"`
	Feature     string `json:"feature"`
	Description string `json:"description"`
	Before      any    `json:"before"`
	After       any    `json:"after"`
	Relation    any    `json:"relation"`
	IpAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
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
	UserAgent               string                   `json:"user_agent"`
	UserResponse            *UserResponse            `json:"user,omitempty" mapstructure:"user"`
	SimpleAuditableResponse *SimpleAuditableResponse `json:"auditable"`
}

type AuditLogPayload struct {
	Before      interface{}
	After       interface{}
	Relation    map[string]interface{}
	Description string
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
	AUDIT_LOG_FEATURE_USER                          = "USER"
	AUDIT_LOG_FEATURE_ROLE                          = "ROLE"
	AUDIT_LOG_FEATURE_MQTT_BROKER                   = "MQTT_BROKER"
	AUDIT_LOG_FEATURE_REGISTER                      = "REGISTER"
	AUDIT_LOG_FEATURE_SYSTEM_SETTING                = "SYSTEM_SETTING"
	AUDIT_LOG_FEATURE_FACILITY                      = "FACILITY"
	AUDIT_LOG_FEATURE_MACHINE                       = "MACHINE"
	AUDIT_LOG_FEATURE_SMTP_SERVER                   = "SMTP_SERVER"
	AUDIT_LOG_FEATURE_MODBUS_SERVER                 = "MODBUS_SERVER"
	AUDIT_LOG_FEATURE_REPORT_DOCUMENT_TEMPLATE      = "REPORT_DOCUMENT_TEMPLATE"
	AUDIT_LOG_FEATURE_CHECK_SHEET_DOCUMENT_TEMPLATE = "CHECK_SHEET_DOCUMENT_TEMPLATE"
	AUDIT_LOG_FEATURE_BREAKDOWN                     = "BREAKDOWN"
	AUDIT_LOG_FEATURE_PARAMETER                     = "PARAMETER"
)

const (
	AUDIT_LOG_ACTOR_SYSTEM = "SYSTEM"
)
