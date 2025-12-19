package parameter_operation

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindBatchById(gormTransaction *gorm.DB, parameterOperationIds []uint64) ([]*entity.ParameterOperation, error)
	CreateBatch(gormTransaction *gorm.DB, parameterOperations []*entity.ParameterOperation) error
	Update(gormTransaction *gorm.DB, parameterOperation *entity.ParameterOperation) error

	DeleteBatchById(gormTransaction *gorm.DB, parameterOperationIds []uint64) error
}
