package model

import "time"

type CreateTelemetryRequest struct {
	ParameterId uint64
	Value       float64
	Timestamp   time.Time
}

type TelemetryReportFilterRequest struct {
	ReportDocumentTemplateId uint64 `json:"report_document_template_id" validate:"required"`
	Interval                 uint   `json:"interval" validate:"required,gte=1"`
	StartDate                string `json:"start_date" validate:"required,datetime"`
	EndDate                  string `json:"end_date" validate:"required,datetime"`
}

type TelemetryGrouped struct {
	Timestamp time.Time               `json:"timestamp"`
	Entries   []*TelemetryReportValue `json:"entries"`
}

type TelemetryReportValue struct {
	Timestamp   time.Time `json:"timestamp"`
	MachineId   uint64    `json:"machine_id"`
	MachineName string    `json:"machine_name"`
	MachineCode string    `json:"machine_code"`
	ParameterId uint64    `json:"parameter_id"`
	Value       *float64  `json:"value"`
}
