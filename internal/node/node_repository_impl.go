package node

import (
	"go-intconnect-api/internal/entity"
	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (nodeRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Node, error) {
	var nodeEntities []entity.Node
	err := gormTransaction.Find(&nodeEntities).Error
	return nodeEntities, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Node, int64, error) {
	var nodeEntities []entity.Node
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("nodename LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.Node{}).
		Preload("NodeGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&nodeEntities).Error
	gormTransaction.Model(&entity.Node{}).Count(&totalItems)
	return nodeEntities, totalItems, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, nodeId uint64) (*entity.Node, error) {
	var nodeEntity entity.Node
	err := gormTransaction.Model(&entity.Node{}).
		Preload("NodeGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", nodeId).Find(&nodeEntity).Error

	return &nodeEntity, err
}

func (nodeRepositoryImpl *RepositoryImpl) FindByName(currencyName string) *entity.Node {
	//TODO implement me
	panic("implement me")
}

func (nodeRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.Node) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (nodeRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, nodeEntity *entity.Node) error {
	return gormTransaction.Model(nodeEntity).Save(nodeEntity).Error
}

func (nodeRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Node{}).Where("id = ?", id).Delete(entity.Node{}).Error
}
