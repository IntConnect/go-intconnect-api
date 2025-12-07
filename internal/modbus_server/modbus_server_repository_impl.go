package modbus_server

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (modbusServerRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.ModbusServer, error) {
	var modbusServerEntities []entity.ModbusServer
	err := gormTransaction.Find(&modbusServerEntities).Error
	return modbusServerEntities, err
}

func (modbusServerRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.ModbusServer, int64, error) {

	var modbusServerEntities []entity.ModbusServer
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.ModbusServer{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("name ILIKE ? OR email ILIKE ? OR name ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&modbusServerEntities).Error; err != nil {
		return nil, 0, err
	}

	return modbusServerEntities, totalItems, nil
}

func (modbusServerRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, modbusServerId uint64) (*entity.ModbusServer, error) {
	var modbusServerEntity entity.ModbusServer
	err := gormTransaction.Model(&entity.ModbusServer{}).
		Where("id = ?", modbusServerId).
		Find(&modbusServerEntity).Error

	return &modbusServerEntity, err
}

func (modbusServerRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.ModbusServer) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (modbusServerRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, modbusServerEntity *entity.ModbusServer) error {
	return gormTransaction.Model(modbusServerEntity).Save(modbusServerEntity).Error
}

func (modbusServerRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.ModbusServer{}).Where("id = ?", id).Delete(&entity.ModbusServer{}).Error
}
