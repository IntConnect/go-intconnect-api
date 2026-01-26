package model

import "mime/multipart"

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
	Parameters        []*ParameterResponse       `json:"parameters" mapstructure:"Parameters"`
	Registers         []*RegisterResponse        `json:"registers" mapstructure:"Registers"`
	MachineDocuments  []*MachineDocumentResponse `json:"machine_documents" mapstructure:"MachineDocuments"`
	DashboardWidget   []*DashboardWidgetResponse `json:"widgets" mapstructure:"DashboardWidgets"`
	Facility          *FacilityResponse          `json:"facility"`
	AuditableResponse *AuditableResponse         `json:"auditable" mapstructure:"-"`
}

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
	ParameterId               *uint64                        `form:"parameter_id" validate:"omitempty,required,exists=parameters;id"`
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

type MachineDashboardWidget struct {
	MachineId          uint64            `json:"-" validate:"required,exists=machines;id"`
	AddedParameterIds  []uint64          `json:"added_parameter_ids"`
	RemoveParameterIds []uint64          `json:"removed_parameter_id"`
	AddedWidgets       []DashboardWidget `json:"added_widgets"`
	EditedWidgets      []DashboardWidget `json:"edited_widgets"`
	RemovedWidgets     []string          `json:"removed_widgets"`
}

type DashboardWidget struct {
	Id     uint64                 `json:"id" `
	Code   string                 `json:"code" validate:"required"`
	Layout map[string]interface{} `json:"layout"`
	Config map[string]interface{} `json:"config"`
}

type DashboardWidgetResponse struct {
	Id     uint64                 `json:"id"`
	Code   string                 `json:"code"`
	Layout map[string]interface{} `json:"layout"`
	Config map[string]interface{} `json:"config"`
}

func (machineResponse *MachineResponse) GetAuditableResponse() *AuditableResponse {
	return machineResponse.AuditableResponse
}

func (machineResponse *MachineResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	machineResponse.AuditableResponse = auditableResponse
}
