package entity

import (
	"encoding/json"
	"go-intconnect-api/internal/trait"

	"gorm.io/gorm"
)

type Facility struct {
	Id            uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string                 `gorm:"column:name"`
	Code          string                 `gorm:"column:code"`
	Description   string                 `gorm:"column:description"`
	Location      string                 `gorm:"column:location"`
	Status        trait.FacilityStatus   `gorm:"column:status"`
	ThumbnailPath string                 `gorm:"column:thumbnail_path"`
	Metadata      map[string]interface{} `gorm:"-"`
	MetadataRaw   []byte                 `gorm:"column:metadata;type:jsonb"`
	Auditable
}

func (facilityEntity *Facility) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(facilityEntity.MetadataRaw) > 0 {
		err = json.Unmarshal(facilityEntity.MetadataRaw, &facilityEntity.Metadata)
	}
	return
}

func (facilityEntity *Facility) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if facilityEntity.Metadata != nil {
		facilityEntity.MetadataRaw, err = json.Marshal(facilityEntity.Metadata)
	}
	return
}

func (facilityEntity Facility) GetAuditable() *Auditable {
	return &facilityEntity.Auditable
}
