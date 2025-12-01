package model

type MqttBrokerResponse struct {
	Id                uint64             `json:"id"`
	HostName          string             `json:"host_name"`
	MqttPort          string             `json:"mqtt_port"`
	WsPort            string             `json:"ws_port"`
	Username          string             `json:"username"`
	Password          string             `json:"password"`
	IsActive          bool               `json:"is_active"`
	AuditableResponse *AuditableResponse `json:"auditable"`
}

type CreateMqttBrokerRequest struct {
	HostName string `json:"host_name" validate:"required,min=3,max=100"`
	MqttPort string `json:"mqtt_port" validate:"required,min=1,max=5"`
	WsPort   string `json:"ws_port" validate:"required,min=1,max=5"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type UpdateMqttBrokerRequest struct {
	Id       uint64 `json:"-" validate:"required,number,gt=0"`
	HostName string `json:"host_name" validate:"required,min=3,max=100"`
	MqttPort string `json:"mqtt_port" validate:"required,min=1,max=5"`
	WsPort   string `json:"ws_port" validate:"required,min=1,max=5"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active" validate:"required"`
}

func (mqttBrokerResponse *MqttBrokerResponse) GetAuditableResponse() *AuditableResponse {
	return mqttBrokerResponse.AuditableResponse
}

func (mqttBrokerResponse *MqttBrokerResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	mqttBrokerResponse.AuditableResponse = auditableResponse
}
