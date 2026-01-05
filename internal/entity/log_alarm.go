package entity

import "time"

type LogAlarm struct {
	Id             uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	ParameterId    uint64     `gorm:"column:parameter_id"`
	Value          float64    `gorm:"column:value"`
	Type           string     `gorm:"column:type"`
	IsActive       bool       `gorm:"column:is_active"`
	AcknowledgedAt *time.Time `gorm:"column:acknowledged_at"`
	FinishedAt     *time.Time `gorm:"column:finished_at"`
	Note           string     `gorm:"column:note"`
	CreatedAt      time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
