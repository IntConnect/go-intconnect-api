package seeders

import (
	"go-intconnect-api/internal/entity"
	"time"

	"gorm.io/gorm"
)

type AlarmLogSeeder struct{}

func (alarmLogSeeder *AlarmLogSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Exec("TRUNCATE TABLE alarm_logs RESTART IDENTITY CASCADE")

	gormDatabase.Model(&entity.AlarmLog{}).Create(&entity.AlarmLog{
		ParameterId:    1,
		AcknowledgedBy: nil,
		Value:          0,
		Type:           "Open",
		IsActive:       true,
		Status:         "",
		Note:           "",
		AcknowledgedAt: nil,
		ResolvedAt:     nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	return nil
}
