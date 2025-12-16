package entity

type Machine struct {
	Id               uint64            `gorm:"column:id;primaryKey;autoIncrement"`
	FacilityId       uint64            `gorm:"column:facility_id;"`
	Name             string            `gorm:"column:name"`
	Code             string            `gorm:"column:code"`
	Description      string            `gorm:"column:description"`
	ModelOffsetX     float32           `gorm:"column:model_offset_x"`
	ModelOffsetY     float32           `gorm:"column:model_offset_y"`
	ModelOffsetZ     float32           `gorm:"column:model_offset_z"`
	ModelScale       float32           `gorm:"column:model_scale"`
	ThumbnailPath    string            `gorm:"column:thumbnail_path"`
	ModelPath        string            `gorm:"column:model_path"`
	Facility         Facility          `gorm:"foreignKey:FacilityId;references:Id"`
	MachineDocuments []MachineDocument `gorm:"foreignKey:MachineId;references:Id"`
	MqttTopic        MqttTopic         `gorm:"foreignKey:MqttTopicId;references:Id"`
	Auditable
}

func (machineEntity Machine) GetAuditable() *Auditable {
	return &machineEntity.Auditable
}
