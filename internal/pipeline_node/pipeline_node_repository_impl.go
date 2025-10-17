package pipeline_node

import (
	"fmt"
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (nodeRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.PipelineNode, error) {
	var nodeEntities []entity.PipelineNode
	err := gormTransaction.Find(&nodeEntities).Error
	fmt.Println(nodeEntities)
	return nodeEntities, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.PipelineNode, int64, error) {
	var nodeEntities []entity.PipelineNode
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("nodename LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.PipelineNode{}).
		Preload("PipelineNodeGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&nodeEntities).Error
	gormTransaction.Model(&entity.PipelineNode{}).Count(&totalItems)
	return nodeEntities, totalItems, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, nodeId uint64) (*entity.PipelineNode, error) {
	var nodeEntity entity.PipelineNode
	err := gormTransaction.Model(&entity.PipelineNode{}).
		Preload("PipelineNodeGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", nodeId).Find(&nodeEntity).Error

	return &nodeEntity, err
}

func (nodeRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineNodeEntity *entity.PipelineNode) error {
	fmt.Println(pipelineNodeEntity)
	return gormTransaction.Omit("Config", "ConfigRaw", "Pipeline", "Node").Create(pipelineNodeEntity).Error

}

func (nodeRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, nodeEntity *entity.PipelineNode) error {
	return gormTransaction.Omit("Config").Create(nodeEntity).Error
}

func (nodeRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.PipelineNode{}).Where("id = ?", id).Delete(entity.PipelineNode{}).Error
}
