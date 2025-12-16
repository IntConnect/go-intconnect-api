package entity

import (
	"go-intconnect-api/internal/trait"
)

type Facility struct {
	Id            uint64               `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string               `gorm:"column:name"`
	Code          string               `gorm:"column:code"`
	Description   string               `gorm:"column:description"`
	Location      string               `gorm:"column:location"`
	Status        trait.FacilityStatus `gorm:"column:status"`
	ThumbnailPath string               `gorm:"column:thumbnail_path"`
	ModelPath     string               `gorm:"column:model_path"`
	PositionX     float64              `gorm:"column:position_x"`
	PositionY     float64              `gorm:"column:position_y"`
	PositionZ     float64              `gorm:"column:position_z"`
	CameraX       float64              `gorm:"column:camera_x;default:0"`
	CameraY       float64              `gorm:"column:camera_y;default:0"`
	CameraZ       float64              `gorm:"column:camera_z;default:0"`
	Auditable
}

func (facilityEntity *Facility) GetAuditable() *Auditable {
	return &facilityEntity.Auditable
}
