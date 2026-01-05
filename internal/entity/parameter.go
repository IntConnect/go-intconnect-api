package entity

import "go-intconnect-api/internal/trait"

type Parameter struct {
	Id                                   uint64                                 `gorm:"column:id;primaryKey;autoIncrement"`
	MqttTopicId                          *uint64                                `gorm:"column:mqtt_topic_id;"`
	Name                                 string                                 `gorm:"column:name"`
	Code                                 string                                 `gorm:"column:code"`
	Unit                                 string                                 `gorm:"column:unit"`
	MinValue                             float64                                `gorm:"column:min_value"`
	MaxValue                             float64                                `gorm:"column:max_value"`
	Description                          string                                 `gorm:"column:description"`
	Category                             trait.ParameterCategory                `gorm:"column:category"`
	PositionX                            float32                                `gorm:"column:position_x"`
	PositionY                            float32                                `gorm:"column:position_y"`
	PositionZ                            float32                                `gorm:"column:position_z"`
	RotationX                            float32                                `gorm:"column:rotation_x"`
	RotationY                            float32                                `gorm:"column:rotation_y"`
	RotationZ                            float32                                `gorm:"column:rotation_z"`
	AbnormalDuration                     int                                    `gorm:"column:abnormal_duration"`
	IsAutomatic                          bool                                   `gorm:"column:is_automatic;"`
	IsDisplay                            bool                                   `gorm:"column:is_display;"`
	IsWatch                              bool                                   `gorm:"column:is_watch;"`
	IsRunningTime                        bool                                   `gorm:"column:is_running_time;"`
	IsFeatured                           bool                                   `gorm:"column:is_featured;"`
	MqttTopic                            *MqttTopic                             `gorm:"foreignKey:MqttTopicId;references:Id"`
	ReportDocumentTemplates              []*ReportDocumentTemplate              `gorm:"many2many:report_document_templates_parameters;joinForeignKey:ParameterID;joinReferences:ReportDocumentTemplateID"`
	ParameterOperations                  []*ParameterOperation                  `gorm:"foreignKey:ParameterId;references:Id"`
	CheckSheetDocumentTemplateParameters []*CheckSheetDocumentTemplateParameter `gorm:"foreignKey:ParameterId;references:Id"`
	Auditable                            `gorm:"embedded"`
}

func (parameterEntity *Parameter) GetAuditable() *Auditable {
	return &parameterEntity.Auditable
}

func (parameterEntity *Parameter) GetId() uint64 {
	return parameterEntity.Id
}
