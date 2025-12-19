package entity

type Parameter struct {
	Id                      uint64                    `gorm:"column:id;primaryKey;autoIncrement"`
	MqttTopicId             *uint64                   `gorm:"column:mqtt_topic_id;"`
	Name                    string                    `gorm:"column:name"`
	Code                    string                    `gorm:"column:code"`
	Unit                    string                    `gorm:"column:unit"`
	MinValue                float32                   `gorm:"column:min_value"`
	MaxValue                float32                   `gorm:"column:max_value"`
	Description             string                    `gorm:"column:description"`
	PositionX               float32                   `gorm:"column:position_x"`
	PositionY               float32                   `gorm:"column:position_y"`
	PositionZ               float32                   `gorm:"column:position_z"`
	RotationX               float32                   `gorm:"column:rotation_x"`
	RotationY               float32                   `gorm:"column:rotation_y"`
	RotationZ               float32                   `gorm:"column:rotation_z"`
	IsDisplay               bool                      `gorm:"column:is_display;default:true"`
	IsAutomatic             bool                      `gorm:"column:is_automatic;default:true"`
	MqttTopic               *MqttTopic                `gorm:"foreignKey:MqttTopicId;references:Id"`
	ReportDocumentTemplates []*ReportDocumentTemplate `gorm:"many2many:report_document_templates_parameters;joinForeignKey:ParameterID;joinReferences:ReportDocumentTemplateID"`
	ParameterOperations     []*ParameterOperation     `gorm:"foreignKey:ParameterId;references:Id"`
	Auditable               `gorm:"embedded"`
}

func (parameterEntity *Parameter) GetAuditable() *Auditable {
	return &parameterEntity.Auditable
}

func (parameterEntity *Parameter) GetId() uint64 {
	return parameterEntity.Id
}
