package database_connection

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.DatabaseConnection, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.DatabaseConnection, int64, error)
	FindById(gormTransaction *gorm.DB, databaseConnectionId uint64) (*entity.DatabaseConnection, error)
	FindByName(databaseConnectionName string) *entity.DatabaseConnection
	Create(gormTransaction *gorm.DB, databaseConnectionEntity *entity.DatabaseConnection) error
	Update(gormTransaction *gorm.DB, databaseConnectionEntity *entity.DatabaseConnection) error
	Delete(gormTransaction *gorm.DB, databaseConnectionId uint64) error
	FindAllByIds(gormTransaction *gorm.DB, databaseConnectionIds []uint64) ([]entity.DatabaseConnection, error)
}
