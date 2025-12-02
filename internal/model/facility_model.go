package model

import "mime/multipart"

type CreateFacilityRequest struct {
	Name            string                 `form:"name" validate:"required,min=3,max=100"`
	Code            string                 `form:"code" validate:"required,min=3,max=100"`
	Description     string                 `form:"description" validate:"omitempty,min=3,max=100"`
	Location        string                 `form:"location" validate:"omitempty,min=3,max=100"`
	Metadata        map[string]interface{} `form:"metadata"`
	ThumbnailHeader *multipart.FileHeader  `form:"thumbnail_header" validate:"required,requiredFile,fileExtension=.png .jpg"`
}

type UpdateFacilityRequest struct {
	Id           uint64                 `json:"-"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Location     string                 `json:"location"`
	Status       string                 `json:"status"`
	ThumbnailUrl string                 `json:"thumbnail_url"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type FacilityResponse struct {
	Id                uint64                 `json:"id"`
	Name              string                 `json:"name"`
	Code              string                 `json:"code"`
	Description       string                 `json:"description"`
	Location          string                 `json:"location"`
	Status            string                 `json:"status"`
	ThumbnailUrl      string                 `json:"thumbnail_url"`
	Metadata          map[string]interface{} `json:"metadata"`
	AuditableResponse *AuditableResponse     `json:"auditable"`
}

func (facilityResponse *FacilityResponse) GetAuditableResponse() *AuditableResponse {
	return facilityResponse.AuditableResponse
}

func (facilityResponse *FacilityResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	facilityResponse.AuditableResponse = auditableResponse
}
