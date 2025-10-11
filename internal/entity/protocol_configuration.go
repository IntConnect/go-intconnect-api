package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ProtocolConfiguration struct {
	Id                 uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Name               string                 `gorm:"column:id"`
	Protocol           string                 `gorm:"column:protocol"`
	Description        string                 `gorm:"column:description"`
	SpecificSetting    map[string]interface{} `gorm:"column:-"`
	SpecificSettingRaw []byte                 `gorm:"column:specific_setting"`
	Auditable
}

func (protocolConfigurationEntity *ProtocolConfiguration) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(protocolConfigurationEntity.SpecificSettingRaw) > 0 {
		err = json.Unmarshal(protocolConfigurationEntity.SpecificSettingRaw, &protocolConfigurationEntity.SpecificSetting)
	}
	return
}

func (protocolConfigurationEntity *ProtocolConfiguration) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if protocolConfigurationEntity.SpecificSetting != nil {
		protocolConfigurationEntity.SpecificSettingRaw, err = json.Marshal(protocolConfigurationEntity.SpecificSetting)
	}
	return
}
