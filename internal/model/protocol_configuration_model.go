package model

type ProtocolConfigurationResponse struct {
	Id                uint64             `json:"id"`
	Name              string             `json:"name"`
	Protocol          string             `json:"protocol"`
	Description       string             `json:"description"`
	AuditableResponse *AuditableResponse `json:"auditable_response"`
}

type CreateProtocolConfigurationDto struct {
	Name               string                 `json:"name" validate:"required"`
	Protocol           string                 `json:"protocol" validate:"required"`
	Description        string                 `json:"description"`
	SpecificSettingRaw map[string]interface{} `gorm:"column:specific_setting_raw"`
}

type UpdateProtocolConfigurationDto struct {
	Id                 uint64                 `json:"id" validate:"required"`
	Name               string                 `json:"name" validate:"required"`
	Protocol           string                 `json:"protocol" validate:"required"`
	Description        string                 `json:"description"`
	SpecificSettingRaw map[string]interface{} `gorm:"column:specific_setting_raw"`
}

type DeleteProtocolConfigurationDto struct {
	Id uint64 `json:"id" validate:"required"`
}
