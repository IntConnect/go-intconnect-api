package model

import "mime/multipart"

type SettingSchema struct {
	NewPayload func() interface{}
}

var SystemSettingSchemas = map[string]SettingSchema{
	"DASHBOARD_SETTINGS": {
		NewPayload: func() interface{} {
			return &DashboardSettingPayload{}
		},
	},
	"LISTENER_SETTINGS": {
		NewPayload: func() interface{} {
			return &ListenerSettingPayload{}
		},
	},
}

type ManageSystemSettingRequest struct {
	Key         string                 `form:"key" validate:"required,min=3,max=100"`
	Description string                 `form:"description"`
	Value       map[string]interface{} `form:"value"`
}

type DashboardSettingPayload struct {
	Showing       string                `json:"showing" validate:"required"`
	MachineId     *uint64               `json:"machine_id" validate:"omitempty,gte=1,exists=machines;id"`
	CameraX       float64               `json:"camera_x" mapstructure:"camera_x" validate:"required"`
	CameraY       float64               `json:"camera_y" mapstructure:"camera_y" validate:"required"`
	CameraZ       float64               `json:"camera_z" mapstructure:"camera_z" validate:"required"`
	PinObjectName string                `json:"pin_object_name" mapstructure:"pin_object_name" validate:""`
	ModelFile     *multipart.FileHeader `json:"model_file" mapstructure:"model_file"`
}

type ListenerSettingPayload struct {
	InsertionWorkersAmount *uint64 `mapstructure:"insertion_workers_amount" validate:"required,gte=1"`
	InsertionQueueSize     *uint64 `mapstructure:"insertion_queue_size" validate:"required,gte=1"`
	ParameterRecoveryCount uint64  `mapstructure:"parameter_recovery_count" validate:"required,gte=1"`
	SnapshotTicker         uint64  `mapstructure:"snapshot_ticker" validate:"required,gte=1"`
	SnapshotTickerType     string  `mapstructure:"snapshot_ticker_type" validate:"required"`
}

type SystemSettingResponse struct {
	Id                uint64                 `json:"-"`
	Key               string                 `json:"key"`
	Description       string                 `json:"description"`
	Value             map[string]interface{} `json:"value"`
	AuditableResponse *AuditableResponse     `json:"auditable_response"`
}

type MinimalSystemSettingResponse struct {
	Value map[string]interface{} `json:"value"`
}

func (systemSettingResponse *SystemSettingResponse) GetAuditableResponse() *AuditableResponse {
	return systemSettingResponse.AuditableResponse
}

func (systemSettingResponse *SystemSettingResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	systemSettingResponse.AuditableResponse = auditableResponse
}
