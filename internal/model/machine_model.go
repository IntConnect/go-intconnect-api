package model

type CreateMachineRequest struct {
	FacilityId   uint64  `json:"facility_id"`
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	ModelOffsetX float32 `json:"model_offset_x"`
	ModelOffsetY float32 `json:"model_offset_y"`
	ModelOffsetZ float32 `json:"model_offset_z"`
	ModelScale   float32 `json:"model_scale"`
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
