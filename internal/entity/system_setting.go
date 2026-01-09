package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type SystemSetting struct {
	Id          uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Key         string                 `gorm:"column:key"`
	Description string                 `gorm:"column:description"`
	Value       map[string]interface{} `gorm:"-:all"`
	ValueRaw    []byte                 `gorm:"column:value;type:jsonb"`
	Auditable   Auditable              `gorm:"embedded"`
}

func (systemSettingEntity *SystemSetting) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(systemSettingEntity.ValueRaw) > 0 {
		err = json.Unmarshal(systemSettingEntity.ValueRaw, &systemSettingEntity.Value)
	}
	return
}

func (systemSettingEntity *SystemSetting) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if systemSettingEntity.Value != nil {
		systemSettingEntity.ValueRaw, err = json.Marshal(systemSettingEntity.Value)
	}
	return
}

func (systemSettingEntity *SystemSetting) GetAuditable() *Auditable {
	return &systemSettingEntity.Auditable
}
