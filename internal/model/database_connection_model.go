package model

type DatabaseConnectionResponse struct {
	Id                uint64                            `json:"id"`
	Name              string                            `json:"name"`
	DatabaseType      string                            `json:"database_type"`
	DatabaseName      string                            `json:"database_name"`
	Description       string                            `json:"description"`
	Config            *DatabaseConnectionConfigResponse `json:"config" mapstructure:"-"`
	AuditableResponse *AuditableResponse                `json:"auditable_response" mapstructure:"-"`
}

type CreateDatabaseConnectionRequest struct {
	Name         string                 `json:"name" validate:"required"`
	DatabaseType string                 `json:"database_type" validate:"required"`
	Description  string                 `json:"description"`
	Config       map[string]interface{} `json:"config"`
	IsActive     bool                   `json:"is_active"`
}

type UpdateDatabaseConnectionRequest struct {
	Id           uint64                 `json:"id" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	DatabaseType string                 `json:"database_type" validate:"required"`
	Description  string                 `json:"description"`
	Config       map[string]interface{} `json:"config"`
}

type DeleteDatabaseConnectionRequest struct {
	Id uint64 `json:"id" validate:"required"`
}

type DatabaseConnectionConfigResponse struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}
