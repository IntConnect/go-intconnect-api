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
}

type ManageSystemSettingRequest struct {
	Key         string                 `form:"key" validate:"required,min=3,max=100"`
	Description string                 `form:"description"`
	Value       map[string]interface{} `form:"value"`
}

type DashboardSettingPayload struct {
	CameraX       float64               `json:"camera_x" mapstructure:"camera_x" validate:"required"`
	CameraY       float64               `json:"camera_y" mapstructure:"camera_y" validate:"required"`
	CameraZ       float64               `json:"camera_z" mapstructure:"camera_z" validate:"required"`
	PinObjectName string                `json:"pin_object_name" mapstructure:"pin_object_name" validate:"required"`
	ModelFile     *multipart.FileHeader `json:"model_file" mapstructure:"model_file" validate:"required"`
}

type SystemSettingResponse struct {
	Id                uint64                 `json:"-"`
	Key               string                 `json:"key"`
	Description       string                 `json:"description"`
	Value             map[string]interface{} `json:"value"`
	AuditableResponse *AuditableResponse     `json:"auditable_response"`
}

func (systemSettingResponse *SystemSettingResponse) GetAuditableResponse() *AuditableResponse {
	return systemSettingResponse.AuditableResponse
}

func (systemSettingResponse *SystemSettingResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	systemSettingResponse.AuditableResponse = auditableResponse
}
