package model

import "database/sql"

type DatabaseConnectionResponse struct {
	Id                uint64                            `json:"id"`
	Name              string                            `json:"name"`
	DatabaseType      string                            `json:"database_type"`
	Schemas           []TableSchema                     `json:"schemas"`
	Config            *DatabaseConnectionConfigResponse `json:"config"`
	AuditableResponse *AuditableResponse                `json:"auditable_response"`
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

type DatabaseTableResponse struct {
	Name string `json:"name"`
}

type TableSchema struct {
	TableName string        `json:"table_name"`
	Columns   []TableColumn `json:"columns"`
}

type TableColumn struct {
	TableName     string         `json:"table_name"`
	ColumnName    string         `json:"column_name"`
	DataType      string         `json:"data_type"`
	IsNullable    string         `json:"is_nullable"`
	ColumnDefault sql.NullString `json:"column_default"`
}
