package facility

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.Facility, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.Facility, int64, error)
	FindById(gormTransaction *gorm.DB, facilityId uint64) (*entity.Facility, error)
	FindByName(facilityName string) *entity.Facility
	Create(gormTransaction *gorm.DB, facilityEntity *entity.Facility) error
	Update(gormTransaction *gorm.DB, facilityEntity *entity.Facility) error
	Delete(gormTransaction *gorm.DB, facilityId uint64) error
	FindBatchById(gormTransaction *gorm.DB, facilityIds []uint64) ([]entity.Facility, error)
}
