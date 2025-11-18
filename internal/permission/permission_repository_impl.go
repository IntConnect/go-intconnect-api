package permission

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (permissionRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Permission, error) {
	var permissionEntities []entity.Permission
	err := gormTransaction.Find(&permissionEntities).Error
	return permissionEntities, err
}

func (permissionRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Permission, int64, error) {
	var permissionEntities []entity.Permission
	var totalItems int64

	if searchQuery != "" {
	}

	// Count total items
	err := gormTransaction.Model(&entity.Permission{}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&permissionEntities).Error
	gormTransaction.Model(&entity.Permission{}).Count(&totalItems)
	return permissionEntities, totalItems, err
}

func (permissionRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, permissionId uint64) (*entity.Permission, error) {
	var permissionEntity entity.Permission
	err := gormTransaction.Model(&entity.Permission{}).Where("id = ?", permissionId).Find(&permissionEntity).Error

	return &permissionEntity, err
}
