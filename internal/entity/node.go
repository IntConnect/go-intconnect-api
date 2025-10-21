package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Node struct {
	Id               uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Type             string                 `gorm:"column:type"`
	Name             string                 `gorm:"column:name"`
	Label            string                 `gorm:"column:label"`
	Description      string                 `gorm:"column:description"`
	HelpText         string                 `gorm:"column:help_text"`
	Color            string                 `gorm:"column:color"`
	Icon             string                 `gorm:"column:icon"`
	ComponentName    string                 `gorm:"column:component_name"`
	DefaultConfig    map[string]interface{} `gorm:"-:all"`
	DefaultConfigRaw []byte                 `gorm:"column:default_config;type:jsonb"`
	PipelineNode     []PipelineNode         `gorm:"foreignKey:NodeId;references:Id;"`
	Auditable
}

func (nodeEntity *Node) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(nodeEntity.DefaultConfigRaw) > 0 {
		err = json.Unmarshal(nodeEntity.DefaultConfigRaw, &nodeEntity.DefaultConfig)
	}
	return
}

func (nodeEntity *Node) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if nodeEntity.DefaultConfig != nil {
		nodeEntity.DefaultConfigRaw, err = json.Marshal(nodeEntity.DefaultConfig)
	}
	return
}
