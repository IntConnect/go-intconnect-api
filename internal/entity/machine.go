package entity

type Machine struct {
	Id               uint64             `gorm:"column:id;primaryKey;autoIncrement"`
	FacilityId       uint64             `gorm:"column:facility_id;"`
	Name             string             `gorm:"column:name"`
	Code             string             `gorm:"column:code"`
	Description      string             `gorm:"column:description"`
	CameraX          float64            `gorm:"column:camera_x;"`
	CameraY          float64            `gorm:"column:camera_y;"`
	CameraZ          float64            `gorm:"column:camera_z;"`
	ThumbnailPath    string             `gorm:"column:thumbnail_path"`
	ModelPath        string             `gorm:"column:model_path"`
	Facility         Facility           `gorm:"foreignKey:FacilityId;references:Id"`
	MqttTopic        *MqttTopic         `gorm:"foreignKey:MachineId;"`
	MachineDocuments []*MachineDocument `gorm:"foreignKey:MachineId;references:Id"`
	DashboardWidgets []*DashboardWidget `gorm:"foreignKey:MachineId;references:Id"`
	Auditable        Auditable          `gorm:"embedded"`
}

func (machineEntity *Machine) GetAuditable() *Auditable {
	return &machineEntity.Auditable
}
