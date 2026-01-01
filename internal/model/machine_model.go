package model

import "mime/multipart"

type CreateMachineRequest struct {
	FacilityId       uint64                         `form:"facility_id" validate:"required,exists=facilities;id"`
	Name             string                         `form:"name" validate:"required,min=3,max=100,unique=machines;name"`
	Code             string                         `form:"code" validate:"required,min=3,max=100,unique=machines;code"`
	Description      string                         `form:"description"`
	CameraX          *float64                       `form:"camera_x" validate:"required"`
	CameraY          *float64                       `form:"camera_y" validate:"required"`
	CameraZ          *float64                       `form:"camera_z" validate:"required"`
	Model            *multipart.FileHeader          `form:"model" validate:"required,requiredFile,fileExtension=.glb"`
	Thumbnail        *multipart.FileHeader          `form:"thumbnail" validate:"required,requiredFile,fileExtension=.png .jpg"`
	MachineDocuments []CreateMachineDocumentRequest `form:"machine_documents" validate:"dive"`
}

type UpdateMachineRequest struct {
	Id                        uint64                         `form:"-" validate:"required,exists=machines;id"`
	FacilityId                uint64                         `form:"facility_id" validate:"required,exists=facilities;id"`
	Name                      string                         `form:"name" validate:"required,min=3,max=100,unique=machines;name;Id"`
	Code                      string                         `form:"code" validate:"required,min=3,max=100,unique=machines;code;Id"`
	Description               string                         `form:"description"`
	CameraX                   *float64                       `form:"camera_x" validate:"required"`
	CameraY                   *float64                       `form:"camera_y" validate:"required"`
	CameraZ                   *float64                       `form:"camera_z" validate:"required"`
	Model                     *multipart.FileHeader          `form:"model" validate:"omitempty,fileExtension=.glb"`
	Thumbnail                 *multipart.FileHeader          `form:"thumbnail" validate:"omitempty,fileExtension=.png .jpg"`
	MachineDocuments          []CreateMachineDocumentRequest `form:"machine_documents" validate:"dive"`
	DeletedMachineDocumentIds []uint64                       `form:"deleted_machine_document_ids" validate:"omitempty,dive"`
}

type MachineResponse struct {
	Id                uint64                     `json:"id"`
	FacilityId        uint64                     `json:"facility_id"`
	Name              string                     `json:"name"`
	Code              string                     `json:"code"`
	Description       string                     `json:"description"`
	CameraX           float64                    `json:"camera_x"`
	CameraY           float64                    `json:"camera_y"`
	CameraZ           float64                    `json:"camera_z"`
	ThumbnailPath     string                     `json:"thumbnail_path"`
	ModelPath         string                     `json:"model_path"`
	MqttTopic         *MqttTopicResponse         `json:"mqtt_topic" mapstructure:"MqttTopic"`
	MachineDocuments  []*MachineDocumentResponse `json:"machine_documents" mapstructure:"MachineDocuments"`
	DashboardWidget   []*DashboardWidget         `json:"widgets" mapstructure:"DashboardWidget"`
	Facility          *FacilityResponse          `json:"facility" mapstructure:"facility"`
	AuditableResponse *AuditableResponse         `json:"auditable" mapstructure:"auditable"`
}

type MachineDashboardWidget struct {
	MachineId       uint64            `json:"-" validate:"required,exists=machines;id"`
	DashboardWidget []DashboardWidget `json:"dashboard_widgets" validate:"required,dive"`
}

type DashboardWidget struct {
	Code   string                 `json:"code" validate:"required"`
	Layout map[string]interface{} `json:"layout"`
	Config map[string]interface{} `json:"config"`
}

func (machineResponse *MachineResponse) GetAuditableResponse() *AuditableResponse {
	return machineResponse.AuditableResponse
}

func (machineResponse *MachineResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	machineResponse.AuditableResponse = auditableResponse
}
