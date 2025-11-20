package mqtt_broker

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.MqttBroker, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.MqttBroker, int64, error)
	FindById(gormTransaction *gorm.DB, mqttBrokerId uint64) (*entity.MqttBroker, error)
	Create(gormTransaction *gorm.DB, mqttBrokerEntity *entity.MqttBroker) error
	Update(gormTransaction *gorm.DB, mqttBrokerEntity *entity.MqttBroker) error
	Delete(gormTransaction *gorm.DB, mqttBrokerId uint64) error
}
