package breakdown

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Breakdown, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.Breakdown, int64, error)
	FindById(gormTransaction *gorm.DB, breakdownId uint64) (*entity.Breakdown, error)
	Create(gormTransaction *gorm.DB, breakdownEntity *entity.Breakdown) error
	Update(gormTransaction *gorm.DB, breakdownEntity *entity.Breakdown) error
	Delete(gormTransaction *gorm.DB, breakdownId uint64) error
}
