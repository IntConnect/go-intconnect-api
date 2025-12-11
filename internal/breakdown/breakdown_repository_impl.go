package breakdown

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (breakdownRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Breakdown, error) {
	var breakdownEntities []entity.Breakdown
	err := gormTransaction.Find(&breakdownEntities).Error
	return breakdownEntities, err
}

func (breakdownRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.Breakdown, int64, error) {

	var breakdownEntities []*entity.Breakdown
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.Breakdown{})

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
		Find(&breakdownEntities).Error; err != nil {
		return nil, 0, err
	}

	return breakdownEntities, totalItems, nil
}

func (breakdownRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, breakdownId uint64) (*entity.Breakdown, error) {
	var breakdownEntity entity.Breakdown
	err := gormTransaction.Model(&entity.Breakdown{}).
		Where("id = ?", breakdownId).
		First(&breakdownEntity).Error

	return &breakdownEntity, err
}

func (breakdownRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.Breakdown) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (breakdownRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, breakdownEntity *entity.Breakdown) error {
	return gormTransaction.Model(breakdownEntity).Save(breakdownEntity).Error
}

func (breakdownRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Breakdown{}).Where("id = ?", id).Delete(&entity.Breakdown{}).Error
}
