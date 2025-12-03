package model

import "mime/multipart"

type CreateFacilityRequest struct {
	Name        string                 `form:"name" validate:"required,min=3,max=100"`
	Code        string                 `form:"code" validate:"required,min=3,max=100,unique=facilities;code"`
	Description string                 `form:"description" validate:"omitempty,min=3,max=100"`
	Location    string                 `form:"location" validate:"omitempty,min=3,max=100"`
	Metadata    map[string]interface{} `form:"metadata"`
	Thumbnail   *multipart.FileHeader  `form:"thumbnail" validate:"required,requiredFile,fileExtension=.png .jpg"`
}

type UpdateFacilityRequest struct {
	Id          uint64                 `json:"-" validate:"required,gt=1"`
	Name        string                 `form:"name" validate:"required,min=3,max=100"`
	Code        string                 `form:"code" validate:"required,min=3,max=100,unique=facilities;code;Id"`
	Description string                 `form:"description" validate:"omitempty,min=3,max=100"`
	Location    string                 `form:"location" validate:"omitempty,min=3,max=100"`
	Metadata    map[string]interface{} `form:"metadata"`
	Thumbnail   *multipart.FileHeader  `form:"thumbnail" validate:"omitempty,fileExtension=.png .jpg"`
}

type FacilityResponse struct {
	Id                uint64                 `json:"id"`
	Name              string                 `json:"name"`
	Code              string                 `json:"code"`
	Description       string                 `json:"description"`
	Location          string                 `json:"location"`
	Status            string                 `json:"status"`
	ThumbnailPath     string                 `json:"thumbnail_path"`
	Metadata          map[string]interface{} `json:"metadata"`
	AuditableResponse *AuditableResponse     `json:"auditable"`
}

func (facilityResponse *FacilityResponse) GetAuditableResponse() *AuditableResponse {
	return facilityResponse.AuditableResponse
}

func (facilityResponse *FacilityResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	facilityResponse.AuditableResponse = auditableResponse
}
