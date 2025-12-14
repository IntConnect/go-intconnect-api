package system_setting

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.SystemSetting, error)
	FindByKey(gormTransaction *gorm.DB, systemSettingKey string) (*entity.SystemSetting, error)
	Manage(gormTransaction *gorm.DB, systemSettingEntity *entity.SystemSetting) error
}
