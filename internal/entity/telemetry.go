package entity

import "time"

type Telemetry struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ParameterId uint64    `gorm:"column:parameter_id;"`
	Value       float64   `gorm:"column:value"`
	Timestamp   time.Time `gorm:"column:timestamp"`
}

type TelemetryQuery struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ParameterId uint64    `gorm:"column:parameter_id;"`
	LastValue   *float64  `gorm:"column:last_value"`
	Bucket      time.Time `gorm:"column:bucket"`
}
