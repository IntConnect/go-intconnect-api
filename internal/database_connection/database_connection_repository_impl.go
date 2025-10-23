package database_connection

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.DatabaseConnection, error) {
	var databaseConnectionEntities []entity.DatabaseConnection
	err := gormTransaction.Find(&databaseConnectionEntities).Error
	return databaseConnectionEntities, err
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) FindAllByIds(gormTransaction *gorm.DB, databaseConnectionIds []uint64) ([]entity.DatabaseConnection, error) {
	var databaseConnectionEntities []entity.DatabaseConnection
	err := gormTransaction.
		Where("id IN ?", databaseConnectionIds).
		Find(&databaseConnectionEntities).Error
	return databaseConnectionEntities, err
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.DatabaseConnection, int64, error) {
	var databaseConnectionEntities []entity.DatabaseConnection
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("name LIKE ? OR protocol LIKE ?  OR description = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.DatabaseConnection{}).
		Preload("DatabaseConnectionGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&databaseConnectionEntities).Error
	gormTransaction.Model(&entity.DatabaseConnection{}).Count(&totalItems)
	return databaseConnectionEntities, totalItems, err
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, databaseConnectionId uint64) (*entity.DatabaseConnection, error) {
	var databaseConnectionEntity entity.DatabaseConnection
	err := gormTransaction.Where("id = ?", databaseConnectionId).Find(&databaseConnectionEntity).Error
	return &databaseConnectionEntity, err
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) FindByName(databaseConnectionName string) *entity.DatabaseConnection {
	//TODO implement me
	panic("implement me")
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, databaseConnectionEntity *entity.DatabaseConnection) error {
	return gormTransaction.Create(databaseConnectionEntity).Error

}

func (databaseConnectionRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, databaseConnectionEntity *entity.DatabaseConnection) error {
	return gormTransaction.Model(databaseConnectionEntity).Save(databaseConnectionEntity).Error
}

func (databaseConnectionRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.DatabaseConnection{}).Where("id = ?", id).Delete(entity.DatabaseConnection{}).Error
}
