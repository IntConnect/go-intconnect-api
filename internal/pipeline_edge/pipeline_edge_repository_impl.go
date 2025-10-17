package pipeline_edge

import (
	"fmt"
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (nodeRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.PipelineEdge, error) {
	var nodeEntities []entity.PipelineEdge
	err := gormTransaction.Find(&nodeEntities).Error
	fmt.Println(nodeEntities)
	return nodeEntities, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.PipelineEdge, int64, error) {
	var nodeEntities []entity.PipelineEdge
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("nodename LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.PipelineEdge{}).
		Preload("PipelineEdgeGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&nodeEntities).Error
	gormTransaction.Model(&entity.PipelineEdge{}).Count(&totalItems)
	return nodeEntities, totalItems, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, nodeId uint64) (*entity.PipelineEdge, error) {
	var nodeEntity entity.PipelineEdge
	err := gormTransaction.Model(&entity.PipelineEdge{}).
		Preload("PipelineEdgeGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", nodeId).Find(&nodeEntity).Error

	return &nodeEntity, err
}

func (nodeRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEdgeEntity *entity.PipelineEdge) error {
	return gormTransaction.Create(pipelineEdgeEntity).Error

}

func (nodeRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, nodeEntity *entity.PipelineEdge) error {
	return gormTransaction.Model(nodeEntity).Save(nodeEntity).Error
}

func (nodeRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.PipelineEdge{}).Where("id = ?", id).Delete(entity.PipelineEdge{}).Error
}
