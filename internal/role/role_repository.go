package role

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Role, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Role, int64, error)
	FindById(gormTransaction *gorm.DB, roleId uint64) (*entity.Role, error)
	Create(gormTransaction *gorm.DB, roleEntity *entity.Role) error
	Update(gormTransaction *gorm.DB, roleEntity *entity.Role) error
	Delete(gormTransaction *gorm.DB, roleId uint64) error
}
