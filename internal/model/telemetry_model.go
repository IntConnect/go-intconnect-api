package model

import "time"

type CreateTelemetryRequest struct {
	ParameterId uint64
	Value       float64
	Timestamp   time.Time
}
