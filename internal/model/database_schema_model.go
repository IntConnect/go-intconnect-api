package model

type CreateDatabaseSchemaRequest struct {
	DatabaseConnectionId uint64            `json:"database_connection_id"`
	TableName            string            `json:"table_name"`
	Columns              []DatabaseColumns `json:"columns"`
}

type DatabaseColumns struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Length        string `json:"length"`
	Nullable      bool   `json:"nullable"`
	Default       string `json:"default"`
	Primary       bool   `json:"primary"`
	AutoIncrement bool   `json:"auto_increment"`
}
