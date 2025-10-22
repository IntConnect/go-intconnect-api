package protocol_configuration

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.ProtocolConfiguration, error) {
	var protocolConfigurationEntities []entity.ProtocolConfiguration
	err := gormTransaction.Find(&protocolConfigurationEntities).Error
	return protocolConfigurationEntities, err
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) FindAllByIds(gormTransaction *gorm.DB, protocolConfigurationIds []uint64) ([]entity.ProtocolConfiguration, error) {
	var protocolConfigurationEntities []entity.ProtocolConfiguration
	err := gormTransaction.
		Where("id IN ?", protocolConfigurationIds).
		Find(&protocolConfigurationEntities).Error
	return protocolConfigurationEntities, err
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.ProtocolConfiguration, int64, error) {
	var protocolConfigurationEntities []entity.ProtocolConfiguration
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("name LIKE ? OR protocol LIKE ?  OR description = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.ProtocolConfiguration{}).
		Preload("ProtocolConfigurationGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&protocolConfigurationEntities).Error
	gormTransaction.Model(&entity.ProtocolConfiguration{}).Count(&totalItems)
	return protocolConfigurationEntities, totalItems, err
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, protocolConfigurationId uint64) (*entity.ProtocolConfiguration, error) {
	var protocolConfigurationEntity entity.ProtocolConfiguration
	err := gormTransaction.Where("id = ?", protocolConfigurationId).Find(&protocolConfigurationEntity).Error
	return &protocolConfigurationEntity, err
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) FindByName(currencyName string) *entity.ProtocolConfiguration {
	//TODO implement me
	panic("implement me")
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.ProtocolConfiguration) error {
	return gormTransaction.Create(currencyEntity).Error

}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, protocolConfigurationEntity *entity.ProtocolConfiguration) error {
	return gormTransaction.Model(protocolConfigurationEntity).Save(protocolConfigurationEntity).Error
}

func (protocolConfigurationRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.ProtocolConfiguration{}).Where("id = ?", id).Delete(entity.ProtocolConfiguration{}).Error
}
