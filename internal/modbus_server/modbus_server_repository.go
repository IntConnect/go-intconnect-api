package modbus_server

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.ModbusServer, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.ModbusServer, int64, error)
	FindById(gormTransaction *gorm.DB, modbusServerId uint64) (*entity.ModbusServer, error)
	Create(gormTransaction *gorm.DB, modbusServerEntity *entity.ModbusServer) error
	Update(gormTransaction *gorm.DB, modbusServerEntity *entity.ModbusServer) error
	Delete(gormTransaction *gorm.DB, modbusServerId uint64) error
}
