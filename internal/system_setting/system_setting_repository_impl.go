package system_setting

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
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

func (systemSettingRepositoryImpl *RepositoryImpl) Manage(gormTransaction *gorm.DB, pipelineEntity *entity.SystemSetting) error {
	return gormTransaction.Model(pipelineEntity).
		Save(pipelineEntity).Error

}
