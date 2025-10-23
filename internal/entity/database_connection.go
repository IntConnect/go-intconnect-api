package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type DatabaseConnection struct {
	Id           uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string                 `gorm:"column:name"`
	DatabaseType string                 `gorm:"column:database_type"`
	DatabaseName string                 `gorm:"column:database_name"`
	Description  string                 `gorm:"column:description"`
	Config       map[string]interface{} `gorm:"-:all"`
	ConfigRaw    []byte                 `gorm:"column:config;type:jsonb"`

	Auditable
}

func (DatabaseConnection *DatabaseConnection) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(DatabaseConnection.ConfigRaw) > 0 {
		err = json.Unmarshal(DatabaseConnection.ConfigRaw, &DatabaseConnection.Config)
	}
	return
}

func (DatabaseConnection *DatabaseConnection) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if DatabaseConnection.Config != nil {
		DatabaseConnection.ConfigRaw, err = json.Marshal(DatabaseConnection.Config)
	}
	return
}
