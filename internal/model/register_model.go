package model

type RegisterResponse struct {
	Id                   uint64                `json:"id"`
	MachineId            uint64                `json:"machine_id"`
	ModbusServerId       uint64                `json:"modbus_server_id"`
	MemoryLocation       string                `json:"memory_location"`
	Name                 string                `json:"name"`
	Description          string                `json:"description"`
	DataType             string                `json:"data_type"`
	MachineResponse      *MachineResponse      `json:"machine" mapstructure:"Machine"`
	ModbusServerResponse *ModbusServerResponse `json:"modbus_server" mapstructure:"ModbusServer"`
	AuditableResponse    *AuditableResponse    `json:"auditable"`
}

type CreateRegisterRequest struct {
	MachineId      uint64 `json:"machine_id" validate:"required,gte=1,exists=machines;id"`
	ModbusServerId uint64 `json:"modbus_server_id" validate:"required,gte=1,exists=modbus_servers;id"`
	MemoryLocation string `json:"memory_location" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
	DataType       string `json:"data_type" validate:"required"`
}

type UpdateRegisterRequest struct {
	Id             uint64 `json:"-" validate:"required,number,gte=1"`
	MachineId      uint64 `json:"machine_id" validate:"required,gte=1,exists=machines;id"`
	ModbusServerId uint64 `json:"modbus_server_id" validate:"required,gte=1,exists=modbus_servers;id"`
	MemoryLocation string `json:"memory_location" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
	DataType       string `json:"data_type" validate:"required"`
}

type RegisterDependency struct {
	MachineResponses      []MachineResponse      `json:"machines"`
	ModbusServerResponses []ModbusServerResponse `json:"modbus_servers"`
}

func (smtpServerResponse *RegisterResponse) GetAuditableResponse() *AuditableResponse {
	return smtpServerResponse.AuditableResponse
}

func (smtpServerResponse *RegisterResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	smtpServerResponse.AuditableResponse = auditableResponse
}
