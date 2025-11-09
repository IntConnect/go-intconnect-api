package utils

import (
	"fmt"
	"go-intconnect-api/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDynamicDatabaseConnection(databaseConnectionResponse *model.DatabaseConnectionResponse) (*gorm.DB, error) {
	var dsn string
	var gormDialector gorm.Dialector

	username := databaseConnectionResponse.Config.Username
	password := databaseConnectionResponse.Config.Password
	host := databaseConnectionResponse.Config.Host
	port := databaseConnectionResponse.Config.Port
	databaseName := databaseConnectionResponse.Config.Database
	databaseDriver := databaseConnectionResponse.DatabaseType
	switch databaseDriver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username, password, host, port, databaseName)
		gormDialector = mysql.Open(dsn)
	case "postgresql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s databasename=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			host, username, password, databaseName, port)
		fmt.Println(dsn)
		gormDialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", databaseDriver)
	}

	database, err := gorm.Open(gormDialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return database, nil
}
