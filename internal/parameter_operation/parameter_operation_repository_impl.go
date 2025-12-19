package parameter_operation

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (parameterOperationRepository *RepositoryImpl) DeleteBatchById(gormTransaction *gorm.DB, parameterOperationIds []uint64) error {
	return gormTransaction.Model(&entity.ParameterOperation{}).Where("id IN ?", parameterOperationIds).Delete(&entity.ParameterOperation{}).Error
}

func (parameterOperationRepository *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, parameterOperations []*entity.ParameterOperation) error {
	err := gormTransaction.Create(&parameterOperations).Error
	return err
}

func (parameterOperationRepository *RepositoryImpl) Update(gormTransaction *gorm.DB, parameterOperationEntity *entity.ParameterOperation) error {
	return gormTransaction.Model(parameterOperationEntity).Save(parameterOperationEntity).Error
}

func (parameterOperationRepository *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, parameterOperationIds []uint64) ([]*entity.ParameterOperation, error) {
	var parameterOperationEntities []*entity.ParameterOperation
	err := gormTransaction.Where("id IN ?", parameterOperationIds).Find(&parameterOperationEntities).Error
	return parameterOperationEntities, err
}
