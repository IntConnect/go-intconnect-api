package mqtt_broker

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (mqttBrokerRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.MqttBroker, error) {
	var mqttBrokerEntities []entity.MqttBroker
	err := gormTransaction.Find(&mqttBrokerEntities).Error
	return mqttBrokerEntities, err
}

func (mqttBrokerRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.MqttBroker, int64, error) {

	var mqttBrokerEntities []entity.MqttBroker
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.MqttBroker{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("host_name ILIKE ? OR mqtt_port ILIKE ? OR ws_port ILIKE ? OR username ILIKE ? OR password ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
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
		Find(&mqttBrokerEntities).Error; err != nil {
		return nil, 0, err
	}

	return mqttBrokerEntities, totalItems, nil
}

func (mqttBrokerRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, mqttBrokerId uint64) (*entity.MqttBroker, error) {
	var mqttBrokerEntity entity.MqttBroker
	err := gormTransaction.Model(&entity.MqttBroker{}).
		Where("id = ?", mqttBrokerId).Find(&mqttBrokerEntity).Error

	return &mqttBrokerEntity, err
}

func (mqttBrokerRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.MqttBroker) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (mqttBrokerRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, mqttBrokerEntity *entity.MqttBroker) error {
	return gormTransaction.Model(mqttBrokerEntity).Save(mqttBrokerEntity).Error
}

func (mqttBrokerRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.MqttBroker{}).Where("id = ?", id).Delete(&entity.MqttBroker{}).Error
}
