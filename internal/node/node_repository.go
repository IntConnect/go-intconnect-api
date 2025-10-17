package node

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Node, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Node, int64, error)
	FindById(gormTransaction *gorm.DB, nodeId uint64) (*entity.Node, error)
	FindByName(nodeName string) *entity.Node
	Create(gormTransaction *gorm.DB, nodeEntity *entity.Node) error
	Update(gormTransaction *gorm.DB, nodeEntity *entity.Node) error
	Delete(gormTransaction *gorm.DB, nodeId uint64) error
	FindBatchById(gormTransaction *gorm.DB, nodeIds []uint64) ([]entity.Node, error)
}
