package model

import "mime/multipart"

type CreateMachineRequest struct {
	FacilityId       uint64                         `form:"facility_id" validate:"required,exists=facilities;id"`
	Name             string                         `form:"name"`
	Code             string                         `form:"code"`
	Description      string                         `form:"description"`
	ModelOffsetX     float32                        `form:"model_offset_x"`
	ModelOffsetY     float32                        `form:"model_offset_y"`
	ModelOffsetZ     float32                        `form:"model_offset_z"`
	ModelScale       float32                        `form:"model_scale"`
	Model            *multipart.FileHeader          `form:"model" validate:"required,requiredFile,fileExtension=.glb"`
	Thumbnail        *multipart.FileHeader          `form:"thumbnail" validate:"required,requiredFile,fileExtension=.png .jpg"`
	MachineDocuments []CreateMachineDocumentRequest `form:"machine_documents" validate:"dive"`
}

type UpdateMachineRequest struct {
	Id                        uint64                         `form:"-" validate:"required,exists=machines;id"`
	FacilityId                uint64                         `form:"facility_id" validate:"required,exists=facilities;id"`
	Name                      string                         `form:"name"`
	Code                      string                         `form:"code"`
	Description               string                         `form:"description"`
	ModelOffsetX              float32                        `form:"model_offset_x"`
	ModelOffsetY              float32                        `form:"model_offset_y"`
	ModelOffsetZ              float32                        `form:"model_offset_z"`
	ModelScale                float32                        `form:"model_scale"`
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
	ModelOffsetX      float32                    `json:"model_offset_x"`
	ModelOffsetY      float32                    `json:"model_offset_y"`
	ModelOffsetZ      float32                    `json:"model_offset_z"`
	ModelScale        float32                    `json:"model_scale"`
	ThumbnailPath     string                     `json:"thumbnail_path"`
	ModelPath         string                     `json:"model_path"`
	Parameters        []ParameterResponse        `json:"parameters" mapstructure:"Parameters"`
	MachineDocuments  []*MachineDocumentResponse `json:"machine_documents" mapstructure:"MachineDocuments"`
	Facility          *FacilityResponse          `json:"facility" mapstructure:"facility"`
	AuditableResponse *AuditableResponse         `json:"auditable" mapstructure:"auditable"`
}

func (machineResponse *MachineResponse) GetAuditableResponse() *AuditableResponse {
	return machineResponse.AuditableResponse
}

func (machineResponse *MachineResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	machineResponse.AuditableResponse = auditableResponse
}
