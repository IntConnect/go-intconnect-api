package breakdown_resource

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (breakdownResourceRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.BreakdownResource, error) {
	var breakdownResourceEntities []entity.BreakdownResource
	err := gormTransaction.Find(&breakdownResourceEntities).Error
	return breakdownResourceEntities, err
}

func (breakdownResourceRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.BreakdownResource, int64, error) {

	var breakdownResourceEntities []*entity.BreakdownResource
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.BreakdownResource{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("name ILIKE ? OR email ILIKE ? OR name ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&breakdownResourceEntities).Error; err != nil {
		return nil, 0, err
	}

	return breakdownResourceEntities, totalItems, nil
}

func (breakdownResourceRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, breakdownResourceId uint64) (*entity.BreakdownResource, error) {
	var breakdownResourceEntity entity.BreakdownResource
	err := gormTransaction.Model(&entity.BreakdownResource{}).
		Where("id = ?", breakdownResourceId).
		First(&breakdownResourceEntity).Error

	return &breakdownResourceEntity, err
}

func (breakdownResourceRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, breakdownResourceEntity *entity.BreakdownResource) error {
	return gormTransaction.Model(breakdownResourceEntity).Create(breakdownResourceEntity).Error

}

func (breakdownResourceRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, breakdownResourceEntities []entity.BreakdownResource) error {
	return gormTransaction.Create(breakdownResourceEntities).Error

}

func (breakdownResourceRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, breakdownResourceEntity *entity.BreakdownResource) error {
	return gormTransaction.Model(breakdownResourceEntity).Save(breakdownResourceEntity).Error
}

func (breakdownResourceRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.BreakdownResource{}).Where("id = ?", id).Delete(&entity.BreakdownResource{}).Error
}
