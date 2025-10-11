package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type PipelineEdge struct {
	Id           uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	PipelineId   uint64                 `gorm:"column:pipeline_id;"`
	EdgeId       uint64                 `gorm:"column:edge_id;"`
	SourceNodeId uint64                 `gorm:"column:source_node_id;"`
	TargetNodeId uint64                 `gorm:"column:target_node_id;"`
	Data         map[string]interface{} `gorm:"-"`
	DataRaw      []byte                 `gorm:"column:config;type:jsonb"`
	Pipeline     Pipeline               `gorm:"foreignKey:PipelineId;references:Id"`
	Node         Node                   `gorm:"foreignKey:NodeId;references:Id"`
	Auditable
}

func (pipelineEdgeEntity *PipelineEdge) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(pipelineEdgeEntity.DataRaw) > 0 {
		err = json.Unmarshal(pipelineEdgeEntity.DataRaw, &pipelineEdgeEntity.Data)
	}
	return
}

func (pipelineEdgeEntity *PipelineEdge) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if pipelineEdgeEntity.Data != nil {
		pipelineEdgeEntity.DataRaw, err = json.Marshal(pipelineEdgeEntity.Data)
	}
	return
}
