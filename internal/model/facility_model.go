package model

type CreateFacilityRequest struct {
	Name        string                 `json:"name"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Location    string                 `json:"location"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateFacilityRequest struct {
	Id           uint64                 `json:"id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Location     string                 `json:"location"`
	Status       string                 `json:"status"`
	ThumbnailUrl string                 `json:"thumbnail_url"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type DeleteFacilityRequest struct {
	Id uint64 `json:"id"`
}

type FacilityResponse struct {
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
