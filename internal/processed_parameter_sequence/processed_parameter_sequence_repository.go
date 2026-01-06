package processed_parameter_sequence

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.ProcessedParameterSequence, error)
	FindById(gormTransaction *gorm.DB, processedParameterSequenceId uint64) (*entity.ProcessedParameterSequence, error)
	FindBatchById(gormTransaction *gorm.DB, processedParameterSequenceIds []uint64) ([]*entity.ProcessedParameterSequence, error)
	Create(gormTransaction *gorm.DB, processedParameterSequenceEntity *entity.ProcessedParameterSequence) error
	CreateBatch(gormTransaction *gorm.DB, processedParameterSequenceEntity []*entity.ProcessedParameterSequence) error
	Delete(gormTransaction *gorm.DB, processedParameterSequenceId uint64) error
	DeleteBatch(gormTransaction *gorm.DB, processedParameterSequenceIds []uint64) error
}
