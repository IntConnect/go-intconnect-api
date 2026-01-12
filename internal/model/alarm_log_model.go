package model

import "time"

type AlarmLogResponse struct {
	Id                 uint64             `json:"id"`
	ParameterId        uint64             `json:"parameter_id"`
	Value              float64            `json:"value"`
	Type               string             `json:"type"`
	IsActive           bool               `json:"is_active"`
	Status             string             `json:"status"`
	AcknowledgedAt     string             `json:"acknowledged_at"`
	Note               string             `json:"note"`
	AcknowledgedByUser *UserResponse      `json:"acknowledged_by" mapstructure:"AcknowledgedByUser"`
	Parameter          *ParameterResponse `json:"parameter" mapstructure:"Parameter"`
	CreatedAt          string             `json:"created_at"`
	UpdatedAt          string             `json:"updated_at"`
}

type UpdateAlarmLogRequest struct {
	Id   uint64 `json:"-" validate:"required,exists=alarm_logs;id"`
	Note string `json:"note" validate:"required"`
}

type AlarmEventType string

const (
	AlarmOpen     AlarmEventType = "OPEN"
	AlarmResolved AlarmEventType = "RESOLVED"
	AlarmFinished AlarmEventType = "FINISHED"
)

type AlarmEvent struct {
	Id          uint64         `json:"id"`
	Type        AlarmEventType `json:"type"`
	ParameterId uint64         `json:"parameter_id"`
	Value       float64        `json:"value,omitempty"`
	Status      string         `json:"status"`
	Timestamp   time.Time      `json:"timestamp"`
}
