package entity

type Machine struct {
	Id            uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	FacilityId    uint64    `gorm:"column:facility_id;"`
	MqttTopicId   uint64    `gorm:"column:mqtt_topic_id;"`
	Name          string    `gorm:"column:name"`
	Code          string    `gorm:"column:code"`
	Description   string    `gorm:"column:description"`
	ModelOffsetX  float32   `json:"model_offset_x"`
	ModelOffsetY  float32   `json:"model_offset_y"`
	ModelOffsetZ  float32   `json:"model_offset_z"`
	ModelScale    float32   `json:"model_scale"`
	ThumbnailPath string    `json:"thumbnail_path"`
	ModelPath     string    `json:"model_path"`
	Facility      Facility  `gorm:"foreignKey:FacilityId;references:Id"`
	MqttTopic     MqttTopic `gorm:"foreignKey:MqttTopicId;references:Id"`
	Auditable
}

func (machineEntity Machine) GetAuditable() *Auditable {
	return &machineEntity.Auditable
}
