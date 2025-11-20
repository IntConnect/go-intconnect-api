package model

type MqttBrokerResponse struct {
	Id       uint64 `json:"id"`
	HostName string `gorm:"json:host_name"`
	MqttPort int    `gorm:"json:mqtt_port"`
	WsPort   int    `gorm:"json:ws_port"`
	Username string `gorm:"json:username"`
	Password string `gorm:"json:password"`
	IsActive bool   `gorm:"json:is_active"`
	AuditableResponse
}

type CreateMqttBrokerRequest struct {
	HostName string `json:"host_name" validate:"required"`
	MqttPort string `json:"mqtt_port" validate:"required"`
	WsPort   string `json:"ws_port" validate:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type UpdateMqttBrokerRequest struct {
	Id       uint64 `json:"id" validate:"required,number"`
	HostName string `gorm:"json:host_name"`
	MqttPort int    `gorm:"json:mqtt_port"`
	WsPort   int    `gorm:"json:ws_port"`
	Username string `gorm:"json:username"`
	Password string `gorm:"json:password"`
	IsActive bool   `gorm:"json:is_active"`
}
