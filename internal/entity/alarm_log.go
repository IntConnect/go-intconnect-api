package entity

import "time"

type AlarmLog struct {
	Id                 uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	ParameterId        uint64     `gorm:"column:parameter_id"`
	AcknowledgedBy     *uint64    `gorm:"column:acknowledged_by"`
	Value              float64    `gorm:"column:value"`
	Type               string     `gorm:"column:type"`
	IsActive           bool       `gorm:"column:is_active"`
	Status             string     `gorm:"column:status"`
	Note               string     `gorm:"column:note"`
	AcknowledgedAt     *time.Time `gorm:"column:acknowledged_at"`
	ResolvedAt         *time.Time `gorm:"column:resolved_at"`
	AcknowledgedByUser *User      `gorm:"foreignKey:AcknowledgedBy;references:Id"`
	Parameter          *Parameter `gorm:"foreignKey:ParameterId;references:Id"`
	CreatedAt          time.Time  `gorm:"column:created_at;<-:create" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
}
