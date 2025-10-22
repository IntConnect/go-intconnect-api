package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type PipelineNode struct {
	Id                    uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	PipelineId            uint64                 `gorm:"column:pipeline_id;"`
	NodeId                uint64                 `gorm:"column:node_id;"`
	TempId                string                 `gorm:"-"`
	Type                  string                 `gorm:"column:type"`
	Label                 string                 `gorm:"column:label"`
	PositionX             float64                `gorm:"column:position_x"`
	PositionY             float64                `gorm:"column:position_y"`
	Config                map[string]interface{} `gorm:"-:all"` // Ignore in all operations
	ConfigRaw             []byte                 `gorm:"column:config;type:jsonb" mapstructure:"-"`
	Appearance            map[string]interface{} `gorm:"-:all"` // Ignore in all operations
	AppearanceRaw         []byte                 `gorm:"column:appearance;type:jsonb" mapstructure:"-"`
	Pipeline              Pipeline               `gorm:"foreignKey:PipelineId;references:Id"`
	Node                  Node                   `gorm:"foreignKey:NodeId;references:Id"`
	ProtocolConfiguration *ProtocolConfiguration `gorm:"-:all"`
}

func (pipelineNodeEntity *PipelineNode) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(pipelineNodeEntity.ConfigRaw) > 0 {
		err = json.Unmarshal(pipelineNodeEntity.ConfigRaw, &pipelineNodeEntity.Config)
	}
	if len(pipelineNodeEntity.AppearanceRaw) > 0 {
		err = json.Unmarshal(pipelineNodeEntity.AppearanceRaw, &pipelineNodeEntity.Appearance)
	}
	return
}

func (pipelineNodeEntity *PipelineNode) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if pipelineNodeEntity.Config != nil {
		pipelineNodeEntity.ConfigRaw, err = json.Marshal(pipelineNodeEntity.Config)
	}
	if pipelineNodeEntity.Appearance != nil {
		pipelineNodeEntity.AppearanceRaw, err = json.Marshal(pipelineNodeEntity.Appearance)
	}
	return
}
