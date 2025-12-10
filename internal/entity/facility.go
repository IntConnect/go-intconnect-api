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
	Auditable
}

func (facilityEntity Facility) GetAuditable() *Auditable {
	return &facilityEntity.Auditable
}
