package breakdown_resource

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.BreakdownResource, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.BreakdownResource, int64, error)
	FindById(gormTransaction *gorm.DB, breakdownResourceId uint64) (*entity.BreakdownResource, error)
	Create(gormTransaction *gorm.DB, breakdownResourceEntity *entity.BreakdownResource) error
	CreateBatch(gormTransaction *gorm.DB, breakdownResourceEntities []entity.BreakdownResource) error
	Update(gormTransaction *gorm.DB, breakdownResourceEntity *entity.BreakdownResource) error
	Delete(gormTransaction *gorm.DB, breakdownResourceId uint64) error
}
