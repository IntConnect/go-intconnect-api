package pipeline

import (
	"fmt"
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (nodeRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Pipeline, error) {
	var nodeEntities []entity.Pipeline
	err := gormTransaction.Find(&nodeEntities).Error
	fmt.Println(nodeEntities)
	return nodeEntities, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Pipeline, int64, error) {
	var nodeEntities []entity.Pipeline
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("nodename LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.Pipeline{}).
		Preload("PipelineGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&nodeEntities).Error
	gormTransaction.Model(&entity.Pipeline{}).Count(&totalItems)
	return nodeEntities, totalItems, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, nodeId uint64) (*entity.Pipeline, error) {
	var nodeEntity entity.Pipeline
	err := gormTransaction.Model(&entity.Pipeline{}).
		Preload("PipelineGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", nodeId).Find(&nodeEntity).Error

	return &nodeEntity, err
}

func (nodeRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Pipeline) error {
	return gormTransaction.Omit("Config", "PipelineNode", "PipelineEdge").Create(pipelineEntity).Error

}

func (nodeRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, nodeEntity *entity.Pipeline) error {
	return gormTransaction.Model(nodeEntity).Save(nodeEntity).Error
}

func (nodeRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Pipeline{}).Where("id = ?", id).Delete(entity.Pipeline{}).Error
}
