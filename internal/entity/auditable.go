package entity

import (
	"time"

	"gorm.io/gorm"
)

type Auditable struct {
	CreatedAt time.Time      `gorm:"<-:create" json:"created_at"`
	CreatedBy string         `gorm:"<-:create" json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	DeletedBy *string        `json:"deleted_by"`
}

func NewAuditable(subject string) Auditable {
	return Auditable{
		CreatedAt: time.Now(),
		CreatedBy: subject,
		UpdatedAt: time.Now(),
		UpdatedBy: subject,
		DeletedAt: gorm.DeletedAt{},
		DeletedBy: nil,
	}
}

func UpdateAuditable(subject string) Auditable {
	return Auditable{
		UpdatedAt: time.Now(),
		UpdatedBy: subject,
	}
}

func DeleteAuditable(subject string) Auditable {
	return Auditable{
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
		DeletedBy: &subject,
	}
}
