package parameter

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB, parameterFilterRequest *model.ParameterFilterRequest) ([]*entity.Parameter, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.Parameter, int64, error)
	FindBatchById(gormTransaction *gorm.DB, parameterIds []uint64) ([]*entity.Parameter, error)
	FindById(gormTransaction *gorm.DB, parameterId uint64) (*entity.Parameter, error)
	Create(gormTransaction *gorm.DB, parameterEntity *entity.Parameter) error
	CreateBatch(gormTransaction *gorm.DB, parameterEntity []*entity.Parameter) error
	Update(gormTransaction *gorm.DB, parameterEntity *entity.Parameter) error
	UpdateBatch(gormTransaction *gorm.DB, parameterEntities []*entity.Parameter) error
	Delete(gormTransaction *gorm.DB, parameterId uint64) error
	FindBatchByMachineId(gormTransaction *gorm.DB, machineId uint64) ([]*entity.Parameter, error)
}
