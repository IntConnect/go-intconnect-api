package parameter

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (parameterRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Parameter, error) {
	var parameterEntities []entity.Parameter
	err := gormTransaction.Find(&parameterEntities).Error
	return parameterEntities, err
}

func (parameterRepositoryImpl *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, parameterIds []uint64) ([]entity.Parameter, error) {
	var parameterEntities []entity.Parameter
	err := gormTransaction.Where("id IN ?", parameterIds).Find(&parameterEntities).Error
	return parameterEntities, err
}

func (parameterRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Parameter, int64, error) {
	var parameterEntities []entity.Parameter
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("parametername LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.Parameter{}).
		Preload("ParameterGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&parameterEntities).Error
	gormTransaction.Model(&entity.Parameter{}).Count(&totalItems)
	return parameterEntities, totalItems, err
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
