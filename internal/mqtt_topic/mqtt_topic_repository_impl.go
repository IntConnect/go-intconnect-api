package mqtt_topic

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (mqttTopicRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.MqttTopic, error) {
	var mqttTopicEntities []entity.MqttTopic
	err := gormTransaction.Find(&mqttTopicEntities).Error
	return mqttTopicEntities, err
}

func (mqttTopicRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.MqttTopic, int64, error) {

	var mqttTopicEntities []*entity.MqttTopic
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.MqttTopic{})

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
		Preload("MqttBroker").
		Preload("Machine").
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&mqttTopicEntities).Error; err != nil {
		return nil, 0, err
	}

	return mqttTopicEntities, totalItems, nil
}

func (mqttTopicRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, mqttTopicId uint64) (*entity.MqttTopic, error) {
	var mqttTopicEntity entity.MqttTopic
	err := gormTransaction.Model(&entity.MqttTopic{}).
		Where("id = ?", mqttTopicId).
		Find(&mqttTopicEntity).Error

	return &mqttTopicEntity, err
}

func (mqttTopicRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.MqttTopic) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (mqttTopicRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, mqttTopicEntity *entity.MqttTopic) error {
	return gormTransaction.Model(mqttTopicEntity).Save(mqttTopicEntity).Error
}

func (mqttTopicRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.MqttTopic{}).Where("id = ?", id).Delete(entity.MqttTopic{}).Error
}
