package processed_parameter_sequence

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (processedParameterSequenceRepository *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.ProcessedParameterSequence, error) {
	var processedParameterSequenceEntities []entity.ProcessedParameterSequence
	err := gormTransaction.Find(&processedParameterSequenceEntities).Error
	return processedParameterSequenceEntities, err
}

func (processedParameterSequenceRepository *RepositoryImpl) FindById(gormTransaction *gorm.DB, processedParameterSequenceId uint64) (*entity.ProcessedParameterSequence, error) {
	var processedParameterSequenceEntity entity.ProcessedParameterSequence
	err := gormTransaction.Model(&entity.ProcessedParameterSequence{}).Where("id = ?", processedParameterSequenceId).Find(&processedParameterSequenceEntity).Error
	return &processedParameterSequenceEntity, err
}

func (processedParameterSequenceRepository *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, processedParameterSequenceIds []uint64) ([]*entity.ProcessedParameterSequence, error) {
	var processedParameterSequenceEntity []*entity.ProcessedParameterSequence
	err := gormTransaction.Model(&entity.ProcessedParameterSequence{}).Where("id IN ?", processedParameterSequenceIds).Find(&processedParameterSequenceEntity).Error
	return processedParameterSequenceEntity, err
}

func (processedParameterSequenceRepository *RepositoryImpl) Create(gormTransaction *gorm.DB, processedParameterSequenceEntity *entity.ProcessedParameterSequence) error {
	return gormTransaction.Create(processedParameterSequenceEntity).Error
}

func (processedParameterSequenceRepository *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, processedParameterSequenceEntity []*entity.ProcessedParameterSequence) error {
	return gormTransaction.Create(processedParameterSequenceEntity).Error
}

func (processedParameterSequenceRepository *RepositoryImpl) Delete(gormTransaction *gorm.DB, processedParameterSequenceId uint64) error {
	return gormTransaction.Model(entity.ProcessedParameterSequence{}).Where("id = ?", processedParameterSequenceId).Delete(entity.ProcessedParameterSequence{}).Error
}

func (processedParameterSequenceRepository *RepositoryImpl) DeleteBatch(gormTransaction *gorm.DB, processedParameterSequenceIds []uint64) error {
	return gormTransaction.Model(entity.ProcessedParameterSequence{}).Where("id = ?", processedParameterSequenceIds).Delete(entity.ProcessedParameterSequence{}).Error
}
