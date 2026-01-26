package model

type RegisterResponse struct {
	Id                   uint64                `json:"id"`
	MachineId            uint64                `json:"machine_id"`
	ModbusServerId       uint64                `json:"modbus_server_id"`
	MemoryLocation       string                `json:"memory_location"`
	Name                 string                `json:"name"`
	Description          string                `json:"description"`
	DataType             string                `json:"data_type"`
	Unit                 string                `json:"unit"`
	PositionX            float64               `json:"position_x"`
	PositionY            float64               `json:"position_y"`
	PositionZ            float64               `json:"position_z"`
	RotationX            float64               `json:"rotation_x"`
	RotationY            float64               `json:"rotation_y"`
	RotationZ            float64               `json:"rotation_z"`
	MachineResponse      *MachineResponse      `json:"machine" mapstructure:"Machine"`
	ModbusServerResponse *ModbusServerResponse `json:"modbus_server" mapstructure:"ModbusServer"`
	AuditableResponse    *AuditableResponse    `json:"auditable"`
}

type CreateRegisterRequest struct {
	MachineId      uint64  `json:"machine_id" validate:"required,gte=1,exists=machines;id"`
	ModbusServerId uint64  `json:"modbus_server_id" validate:"required,gte=1,exists=modbus_servers;id"`
	MemoryLocation string  `json:"memory_location" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	Description    string  `json:"description" validate:"required"`
	DataType       string  `json:"data_type" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	PositionX      float64 `json:"position_x"`
	PositionY      float64 `json:"position_y"`
	PositionZ      float64 `json:"position_z"`
	RotationX      float64 `json:"rotation_x"`
	RotationY      float64 `json:"rotation_y"`
	RotationZ      float64 `json:"rotation_z"`
}

type UpdateRegisterRequest struct {
	Id             uint64  `json:"-" validate:"required,number,gte=1"`
	MachineId      uint64  `json:"machine_id" validate:"required,gte=1,exists=machines;id"`
	ModbusServerId uint64  `json:"modbus_server_id" validate:"required,gte=1,exists=modbus_servers;id"`
	MemoryLocation string  `json:"memory_location" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	Description    string  `json:"description" validate:"required"`
	DataType       string  `json:"data_type" validate:"required"`
	PositionX      float64 `json:"position_x"`
	PositionY      float64 `json:"position_y"`
	PositionZ      float64 `json:"position_z"`
	RotationX      float64 `json:"rotation_x"`
	RotationY      float64 `json:"rotation_y"`
	RotationZ      float64 `json:"rotation_z"`
}

type UpdateRegisterValueRequest struct {
	Id    uint64  `json:"-" validate:"required,number,exists=registers;id,gte=1"`
	Value float32 `json:"value" validate:"required,number"`
}

type RegisterDependency struct {
	MachineResponses      []MachineResponse      `json:"machines"`
	ModbusServerResponses []ModbusServerResponse `json:"modbus_servers"`
}

func (registerResponse *RegisterResponse) GetAuditableResponse() *AuditableResponse {
	return registerResponse.AuditableResponse
}

func (registerResponse *RegisterResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	registerResponse.AuditableResponse = auditableResponse
}
