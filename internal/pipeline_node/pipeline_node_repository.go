package pipeline_node

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.PipelineNode, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.PipelineNode, int64, error)
	FindById(gormTransaction *gorm.DB, pipelineNodeId uint64) (*entity.PipelineNode, error)
	Create(gormTransaction *gorm.DB, pipelineNodeEntity *entity.PipelineNode) error
	Update(gormTransaction *gorm.DB, pipelineNodeEntity *entity.PipelineNode) error
	Delete(gormTransaction *gorm.DB, pipelineNodeId uint64) error
}
