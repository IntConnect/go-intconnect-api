package role

import (
	"context"
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.Role, error)
	FindAllCache(context context.Context) ([]*entity.Role, error)
	FindRoleCacheById(ctx context.Context, roleId uint64) (*entity.Role, error)
	FindById(gormTransaction *gorm.DB, roleId uint64) (*entity.Role, error)
	Create(gormTransaction *gorm.DB, roleEntity *entity.Role) error
	SetAllCache(ctx context.Context, roles []*entity.Role) error
	SetByIdCache(ctx context.Context, roleId uint64, roleEntity *entity.Role) error
	Update(gormTransaction *gorm.DB, roleEntity *entity.Role) error
	Delete(gormTransaction *gorm.DB, roleId uint64) error
	DeleteAllCache(ctx context.Context) error
	DeleteByIdCache(ctx context.Context, roleId uint64) error
}
