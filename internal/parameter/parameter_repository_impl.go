package parameter

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (parameterRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.Parameter, error) {
	var parameterEntities []*entity.Parameter
	err := gormTransaction.Find(&parameterEntities).Error
	return parameterEntities, err
}

func (parameterRepositoryImpl *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, parameterIds []uint64) ([]*entity.Parameter, error) {
	var parameterEntities []*entity.Parameter
	err := gormTransaction.Where("id IN ?", parameterIds).Find(&parameterEntities).Error
	return parameterEntities, err
}

func (parameterRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.Parameter, int64, error) {

	var parameterEntities []*entity.Parameter
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.Parameter{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("name ILIKE ? OR code ILIKE ? OR unit ILIKE ? OR min_value ILIKE ? OR max_value ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
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
		Find(&parameterEntities).Error; err != nil {
		return nil, 0, err
	}

	return parameterEntities, totalItems, nil
}
func (parameterRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, parameterId uint64) (*entity.Parameter, error) {
	var parameterEntity entity.Parameter
	err := gormTransaction.Model(&entity.Parameter{}).Where("id = ?", parameterId).Find(&parameterEntity).Error

	return &parameterEntity, err
}

func (parameterRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, parameterEntity *entity.Parameter) error {
	return gormTransaction.Model(parameterEntity).Create(parameterEntity).Error
}

func (parameterRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, parameterEntity []*entity.Parameter) error {
	return gormTransaction.Model(parameterEntity).Create(parameterEntity).Error
}

func (parameterRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, parameterEntity *entity.Parameter) error {
	return gormTransaction.Model(parameterEntity).Save(parameterEntity).Error
}

func (parameterRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Parameter{}).Where("id = ?", id).Delete(entity.Parameter{}).Error
}
