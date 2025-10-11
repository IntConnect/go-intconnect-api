package pipeline

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Pipeline, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Pipeline, int64, error)
	FindById(gormTransaction *gorm.DB, pipelineId uint64) (*entity.Pipeline, error)
	Create(gormTransaction *gorm.DB, pipelineEntity *entity.Pipeline) error
	Update(gormTransaction *gorm.DB, pipelineEntity *entity.Pipeline) error
	Delete(gormTransaction *gorm.DB, pipelineId uint64) error
}
