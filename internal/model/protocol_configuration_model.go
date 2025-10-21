package model

type ProtocolConfigurationResponse struct {
	Id                uint64                 `json:"id"`
	Name              string                 `json:"name"`
	Protocol          string                 `json:"protocol"`
	Description       string                 `json:"description"`
	Config            map[string]interface{} `json:"config"`
	AuditableResponse *AuditableResponse     `json:"auditable_response"`
}

type CreateProtocolConfigurationRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Protocol    string                 `json:"protocol" validate:"required"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	IsActive    bool                   `json:"is_active"`
}

type UpdateProtocolConfigurationRequest struct {
	Id          uint64                 `json:"id" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Protocol    string                 `json:"protocol" validate:"required"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

type DeleteProtocolConfigurationRequest struct {
	Id uint64 `json:"id" validate:"required"`
}
