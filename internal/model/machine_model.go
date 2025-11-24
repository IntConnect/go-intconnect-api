package model

import "mime/multipart"

type CreateMachineRequest struct {
	FacilityId   uint64                `form:"facility_id"`
	Name         string                `form:"name"`
	Code         string                `form:"code"`
	Description  string                `form:"description"`
	ModelOffsetX float32               `form:"model_offset_x"`
	ModelOffsetY float32               `form:"model_offset_y"`
	ModelOffsetZ float32               `form:"model_offset_z"`
	ModelScale   float32               `form:"model_scale"`
	ModelHeader  *multipart.FileHeader `form:"model_header" validate:"required,requiredFile,maxSize=2,fileExtension=.glb"`
}

type UpdateMachineRequest struct {
	Id           uint64  `json:"id"`
	FacilityId   uint64  `json:"facility_id"`
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	ModelOffsetX float32 `json:"model_offset_x"`
	ModelOffsetY float32 `json:"model_offset_y"`
	ModelOffsetZ float32 `json:"model_offset_z"`
	ModelScale   float32 `json:"model_scale"`
}

type DeleteMachineRequest struct {
	Id uint64 `json:"id"`
}

type MachineResponse struct {
	Id                 uint64                 `json:"id"`
	Name               string                 `json:"name"`
	Code               string                 `json:"code"`
	Description        string                 `json:"description"`
	Location           string                 `json:"location"`
	Status             string                 `json:"status"`
	ThumbnailUrl       string                 `json:"thumbnail_url"`
	Metadata           map[string]interface{} `json:"metadata"`
	*AuditableResponse `json:"auditable"`
}

func (machineResponse *MachineResponse) GetAuditableResponse() *AuditableResponse {
	return machineResponse.AuditableResponse
}

func (machineResponse *MachineResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	machineResponse.AuditableResponse = auditableResponse
}
