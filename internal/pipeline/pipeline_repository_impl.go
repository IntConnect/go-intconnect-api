package pipeline

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (pipelineRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.Pipeline, error) {
	var pipelineEntities []*entity.Pipeline
	err := gormTransaction.
		Find(&pipelineEntities).Error
	return pipelineEntities, err
}

func (pipelineRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.Pipeline, int64, error) {
	var pipelineEntities []*entity.Pipeline
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.Pipeline{}).
		Preload("PipelineGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&pipelineEntities).Error
	gormTransaction.Model(&entity.Pipeline{}).Count(&totalItems)
	return pipelineEntities, totalItems, err
}

func (pipelineRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, pipelineId uint64) (*entity.Pipeline, error) {
	var pipelineEntity entity.Pipeline
	err := gormTransaction.
		Preload("PipelineNode").
		Preload("PipelineNode.Node").
		Preload("PipelineEdge").
		First(&pipelineEntity, "id = ?", pipelineId).Error

	return &pipelineEntity, err
}

func (pipelineRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Pipeline) error {
	return gormTransaction.Omit("PipelineNode", "PipelineEdge").Create(pipelineEntity).Error

}

func (pipelineRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, pipelineEntity *entity.Pipeline) error {
	return gormTransaction.Omit("PipelineNode", "PipelineEdge").Updates(pipelineEntity).Error
}

func (pipelineRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Pipeline{}).Where("id = ?", id).Delete(entity.Pipeline{}).Error
}
