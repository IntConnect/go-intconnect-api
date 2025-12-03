package seeders

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type MqttBrokerSeeder struct{}

func (mqttBrokerSeeder *MqttBrokerSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Exec("TRUNCATE TABLE mqtt_brokers RESTART IDENTITY CASCADE")

	gormDatabase.Model(&entity.MqttBroker{}).Create(&entity.MqttBroker{
		HostName:  "127.0.0.1",
		MqttPort:  "1883",
		WsPort:    "9001",
		Username:  "",
		Password:  "",
		IsActive:  true,
		Auditable: entity.NewAuditable("System"),
	})
	gormDatabase.Model(&entity.MqttBroker{}).Create(&entity.MqttBroker{
		HostName:  "10.175.16.39",
		MqttPort:  "1883",
		WsPort:    "9001",
		Username:  "",
		Password:  "",
		IsActive:  true,
		Auditable: entity.NewAuditable("System"),
	})
	return nil
}
