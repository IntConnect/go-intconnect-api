package system_setting

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (systemSettingRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.SystemSetting, error) {
	var systemSettingEntities []*entity.SystemSetting
	err := gormTransaction.Find(&systemSettingEntities).Error
	return systemSettingEntities, err
}

func (systemSettingRepositoryImpl *RepositoryImpl) FindByKey(gormTransaction *gorm.DB, systemSettingKey string) (*entity.SystemSetting, error) {
	var systemSettingEntity entity.SystemSetting
	err := gormTransaction.Model(&entity.SystemSetting{}).
		Where("key = ?", systemSettingKey).First(&systemSettingEntity).Error

	return &systemSettingEntity, err
}

func (systemSettingRepositoryImpl *RepositoryImpl) Manage(
	gormTransaction *gorm.DB,
	systemSettingEntity *entity.SystemSetting,
) error {

	return gormTransaction.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}}, // UNIQUE key
		DoUpdates: clause.AssignmentColumns([]string{
			"description",
			"value",
			"updated_at",
			"updated_by",
		}),
	}).Create(systemSettingEntity).Error
}
