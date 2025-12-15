package machine

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Machine, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Machine, int64, error)
	FindById(gormTransaction *gorm.DB, machineId uint64) (*entity.Machine, error)
	FindByFacilityId(gormTransaction *gorm.DB, facilityId uint64) ([]*entity.Machine, error)
	FindBatchById(gormTransaction *gorm.DB, machineIds []uint64) ([]entity.Machine, error)
	Create(gormTransaction *gorm.DB, machineEntity *entity.Machine) error
	Update(gormTransaction *gorm.DB, machineEntity *entity.Machine) error
	Delete(gormTransaction *gorm.DB, machineEntity *entity.Machine) error
}
