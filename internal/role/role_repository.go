package role

import (
	"context"
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Role, error)
	FindAllCache(context context.Context) ([]entity.Role, error)
	FindById(gormTransaction *gorm.DB, roleId uint64) (*entity.Role, error)
	Create(gormTransaction *gorm.DB, roleEntity *entity.Role) error
	SetAll(ctx context.Context, roles []entity.Role) error
	Update(gormTransaction *gorm.DB, roleEntity *entity.Role) error
	Delete(gormTransaction *gorm.DB, roleId uint64) error
	DeleteAll(ctx context.Context) error
}
