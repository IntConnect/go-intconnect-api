package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Pipeline struct {
	ID          uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string                 `gorm:"column:name"`
	Description string                 `gorm:"column:description"`
	Config      map[string]interface{} `gorm:"-"`
	ConfigRaw   []byte                 `gorm:"column:config;type:jsonb"`
	Auditable
}

func (pipelineEntity *Pipeline) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(pipelineEntity.ConfigRaw) > 0 {
		err = json.Unmarshal(pipelineEntity.ConfigRaw, &pipelineEntity.Config)
	}
	return
}

func (pipelineEntity *Pipeline) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if pipelineEntity.Config != nil {
		pipelineEntity.ConfigRaw, err = json.Marshal(pipelineEntity.Config)
	}
	return
}
