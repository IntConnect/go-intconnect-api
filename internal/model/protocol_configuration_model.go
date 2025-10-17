package model

type ProtocolConfigurationResponse struct {
	Id                uint64             `json:"id"`
	Name              string             `json:"name"`
	Protocol          string             `json:"protocol"`
	Description       string             `json:"description"`
	AuditableResponse *AuditableResponse `json:"auditable_response"`
}

type CreateProtocolConfigurationRequest struct {
	Name               string                 `json:"name" validate:"required"`
	Protocol           string                 `json:"protocol" validate:"required"`
	Description        string                 `json:"description"`
	SpecificSettingRaw map[string]interface{} `gorm:"column:specific_setting_raw"`
}

type UpdateProtocolConfigurationRequest struct {
	Id                 uint64                 `json:"id" validate:"required"`
	Name               string                 `json:"name" validate:"required"`
	Protocol           string                 `json:"protocol" validate:"required"`
	Description        string                 `json:"description"`
	SpecificSettingRaw map[string]interface{} `gorm:"column:specific_setting_raw"`
}

type DeleteProtocolConfigurationRequest struct {
	Id uint64 `json:"id" validate:"required"`
}
