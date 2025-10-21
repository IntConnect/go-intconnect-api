package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ProtocolConfiguration struct {
	Id          uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string                 `gorm:"column:name"`
	Protocol    string                 `gorm:"column:protocol"`
	Description string                 `gorm:"column:description"`
	Config      map[string]interface{} `gorm:"-:all"`
	ConfigRaw   []byte                 `gorm:"column:config"`
	IsActive    bool                   `gorm:"column:is_active"`
	Auditable
}

func (protocolConfigurationEntity *ProtocolConfiguration) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(protocolConfigurationEntity.ConfigRaw) > 0 {
		err = json.Unmarshal(protocolConfigurationEntity.ConfigRaw, &protocolConfigurationEntity.Config)
	}
	return
}

func (protocolConfigurationEntity *ProtocolConfiguration) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if protocolConfigurationEntity.Config != nil {
		protocolConfigurationEntity.ConfigRaw, err = json.Marshal(protocolConfigurationEntity.Config)
	}
	return
}
