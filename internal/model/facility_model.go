package model

import "mime/multipart"

type CreateFacilityRequest struct {
	Name        string                `form:"name" validate:"required,min=3,max=100"`
	Code        string                `form:"code" validate:"required,min=3,max=100,unique=facilities;code"`
	Description string                `form:"description" validate:"omitempty,min=3,max=100"`
	Location    string                `form:"location" validate:"omitempty,min=3,max=100"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail" validate:"required,requiredFile,fileExtension=.png .jpg"`
	Model       *multipart.FileHeader `form:"model" validate:"required,requiredFile,fileExtension=.glb"`
	PositionX   *float64              `form:"position_x" validate:"required,number"`
	PositionY   *float64              `form:"position_y" validate:"required,number"`
	PositionZ   *float64              `form:"position_z" validate:"required,number"`
}

type UpdateFacilityRequest struct {
	Id          uint64                `json:"-" validate:"required,gte=1"`
	Name        string                `form:"name" validate:"required,min=3,max=100"`
	Code        string                `form:"code" validate:"required,min=3,max=100,unique=facilities;code;Id"`
	Description string                `form:"description" validate:"omitempty,min=3,max=100"`
	Location    string                `form:"location" validate:"omitempty,min=3,max=100"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail" validate:"omitempty,fileExtension=.png .jpg"`
	Model       *multipart.FileHeader `form:"model" validate:"omitempty,fileExtension=.glb"`
	PositionX   *float64              `form:"position_x" validate:"required,number"`
	PositionY   *float64              `form:"position_y" validate:"required,number"`
	PositionZ   *float64              `form:"position_z" validate:"required,number"`
}

type FacilityResponse struct {
	Id                uint64             `json:"id"`
	Name              string             `json:"name"`
	Code              string             `json:"code"`
	Description       string             `json:"description"`
	Location          string             `json:"location"`
	Status            string             `json:"status"`
	ThumbnailPath     string             `json:"thumbnail_path"`
	ModelPath         string             `json:"model_path"`
	PositionX         float64            `json:"position_x"`
	PositionY         float64            `json:"position_y"`
	PositionZ         float64            `json:"position_z"`
	AuditableResponse *AuditableResponse `json:"auditable"`
}

func (facilityResponse *FacilityResponse) GetAuditableResponse() *AuditableResponse {
	return facilityResponse.AuditableResponse
}

func (facilityResponse *FacilityResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	facilityResponse.AuditableResponse = auditableResponse
}
