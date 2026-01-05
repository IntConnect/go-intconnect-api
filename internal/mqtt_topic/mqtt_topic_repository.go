package mqtt_topic

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.MqttTopic, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.MqttTopic, int64, error)
	FindById(gormTransaction *gorm.DB, mqttTopicId uint64) (*entity.MqttTopic, error)
	Create(gormTransaction *gorm.DB, mqttTopicEntity *entity.MqttTopic) error
	Update(gormTransaction *gorm.DB, mqttTopicEntity *entity.MqttTopic) error
	Delete(gormTransaction *gorm.DB, mqttTopicId uint64) error
}
