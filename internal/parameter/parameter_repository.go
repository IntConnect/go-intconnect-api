package parameter

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Machine, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Machine, int64, error)
	FindById(gormTransaction *gorm.DB, machineId uint64) (*entity.Machine, error)
	FindByName(machineName string) *entity.Machine
	Create(gormTransaction *gorm.DB, machineEntity *entity.Machine) error
	Update(gormTransaction *gorm.DB, machineEntity *entity.Machine) error
	Delete(gormTransaction *gorm.DB, machineId uint64) error
	FindBatchById(gormTransaction *gorm.DB, machineIds []uint64) ([]entity.Machine, error)
}
