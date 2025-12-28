package seeders

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type MqttTopicSeeder struct{}

func (mqttTopicSeeder *MqttTopicSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Exec("TRUNCATE TABLE mqtt_topics RESTART IDENTITY CASCADE")

	gormDatabase.Model(&entity.MqttTopic{}).Create(&entity.MqttTopic{
		MachineId:    1,
		MqttBrokerId: 1,
		Name:         "sensor/payload",
		QoS:          0,
		Auditable:    entity.NewAuditable("System"),
	})

	return nil
}
