package pipeline_edge

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.PipelineEdge, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.PipelineEdge, int64, error)
	FindById(gormTransaction *gorm.DB, pipelineEdgeId uint64) (*entity.PipelineEdge, error)
	Create(gormTransaction *gorm.DB, pipelineEdgeEntity *entity.PipelineEdge) error
	Update(gormTransaction *gorm.DB, pipelineEdgeEntity *entity.PipelineEdge) error
	Delete(gormTransaction *gorm.DB, pipelineEdgeId uint64) error
}
