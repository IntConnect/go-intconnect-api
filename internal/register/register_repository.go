package register

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.Register, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.Register, int64, error)
	FindById(gormTransaction *gorm.DB, registerId uint64) (*entity.Register, error)
	Create(gormTransaction *gorm.DB, registerEntity *entity.Register) error
	Update(gormTransaction *gorm.DB, registerEntity *entity.Register) error
	Delete(gormTransaction *gorm.DB, registerId uint64) error
}
