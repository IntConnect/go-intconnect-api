package model

type ModbusServerResponse struct {
	Id                uint64             `json:"id"`
	IpAddress         string             `json:"ip_address"`
	Port              string             `json:"port"`
	SlaveId           string             `json:"slave_id"`
	Timeout           int                `json:"timeout"`
	IsActive          bool               `json:"is_active"`
	AuditableResponse *AuditableResponse `json:"auditable"`
}

type CreateModbusServerRequest struct {
	IpAddress string `json:"ip_address" validate:"required,min=1,max=255"`
	Port      string `json:"port" validate:"required,min=1,max=255"`
	SlaveId   string `json:"slave_id" validate:"required,min=1,max=255"`
	Timeout   int    `json:"timeout" validate:"required"`
	IsActive  bool   `json:"is_active"`
}

type UpdateModbusServerRequest struct {
	Id        uint64 `json:"-" validate:"required,exists=modbus_servers;id"`
	IpAddress string `json:"ip_address" validate:"required,min=1,max=255"`
	Port      string `json:"port" validate:"required,min=1,max=255"`
	SlaveId   string `json:"slave_id" validate:"required,min=1,max=255"`
	Timeout   int    `json:"timeout" validate:"required,gte=0"`
	IsActive  bool   `json:"is_active"`
}

func (mqttTopicResponse *ModbusServerResponse) GetAuditableResponse() *AuditableResponse {
	return mqttTopicResponse.AuditableResponse
}

func (mqttTopicResponse *ModbusServerResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	mqttTopicResponse.AuditableResponse = auditableResponse
}
