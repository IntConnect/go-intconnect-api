package entity

type Facility struct {
	Id            uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string     `gorm:"column:name"`
	Code          string     `gorm:"column:code"`
	Description   string     `gorm:"column:description"`
	Location      string     `gorm:"column:location"`
	ThumbnailPath string     `gorm:"column:thumbnail_path"`
	ModelPath     string     `gorm:"column:model_path"`
	PositionX     float64    `gorm:"column:position_x"`
	PositionY     float64    `gorm:"column:position_y"`
	PositionZ     float64    `gorm:"column:position_z"`
	CameraX       float64    `gorm:"column:camera_x;"`
	CameraY       float64    `gorm:"column:camera_y;"`
	CameraZ       float64    `gorm:"column:camera_z;"`
	Machine       []*Machine `gorm:"foreignKey:FacilityId;references:Id"`
	Auditable     Auditable  `gorm:"embedded"`
}

func (facilityEntity *Facility) GetAuditable() *Auditable {
	return &facilityEntity.Auditable
}
