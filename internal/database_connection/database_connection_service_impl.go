package database_connection

import (
	"database/sql"
	"fmt"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"go-intconnect-api/utils"
	"log"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	databaseConnectionRepository Repository
	validatorService             validator.Service
	dbConnection                 *gorm.DB
	viperConfig                  *viper.Viper
}

func NewService(databaseConnectionRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		databaseConnectionRepository: databaseConnectionRepository,
		validatorService:             validatorService,
		dbConnection:                 dbConnection,
		viperConfig:                  viperConfig,
	}
}

func (databaseConnectionService *ServiceImpl) FindAll() []*model.DatabaseConnectionResponse {
	var allDatabaseConnection []*model.DatabaseConnectionResponse

	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		databaseConnectionResponse, err := databaseConnectionService.databaseConnectionRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allDatabaseConnection = mapper.MapDatabaseConnectionEntitiesIntoDatabaseConnectionResponses(databaseConnectionResponse)

		for _, dbConn := range allDatabaseConnection {
			dynamicDB, err := utils.NewDynamicDatabaseConnection(dbConn)
			fmt.Println(dynamicDB, err)
			if err != nil {
				log.Printf("❌ Failed to connect to %s: %v", dbConn.Name, err)
				continue
			}

			switch dbConn.DatabaseType {
			case "postgresql":
				rows, err := dynamicDB.Raw(`
					SELECT 
						c.table_name,
						c.column_name,
						c.data_type,
						c.is_nullable,
						c.column_default
					FROM information_schema.columns c
					JOIN information_schema.tables t
					  ON c.table_name = t.table_name
					WHERE t.table_schema = 'public'
					ORDER BY c.table_name, c.ordinal_position;
				`).Rows()
				fmt.Println(rows.Next(), err)
				if err != nil {
					log.Println("❌ Error fetching PostgreSQL schema:", err)
					continue
				}
				defer rows.Close()

				schemas := map[string][]model.TableColumn{}
				for rows.Next() {
					var col model.TableColumn
					if err := rows.Scan(&col.TableName, &col.ColumnName, &col.DataType, &col.IsNullable, &col.ColumnDefault); err != nil {
						log.Println("⚠️ Scan error:", err)
						continue
					}
					fmt.Println(col)
					schemas[col.TableName] = append(schemas[col.TableName], col)
				}

				// Simpan schema ke response
				var allSchemas []model.TableSchema
				for tableName, cols := range schemas {
					allSchemas = append(allSchemas, model.TableSchema{
						TableName: tableName,
						Columns:   cols,
					})
				}
				dbConn.Schemas = allSchemas

			case "mysql":
				tableRows, err := dynamicDB.Raw("SHOW TABLES").Rows()
				if err != nil {
					log.Println("❌ Error fetching MySQL tables:", err)
					continue
				}

				defer tableRows.Close()

				var tables []string
				for tableRows.Next() {
					var t string
					tableRows.Scan(&t)
					tables = append(tables, t)
				}

				var allSchemas []model.TableSchema
				for _, table := range tables {
					colRows, err := dynamicDB.Raw(fmt.Sprintf("SHOW COLUMNS FROM `%s`", table)).Rows()
					if err != nil {
						log.Println("❌ Error fetching columns for table:", table, err)
						continue
					}
					defer colRows.Close()

					var columns []model.TableColumn
					for colRows.Next() {
						var field, colType, nullable, key, defaultValue, extra sql.NullString
						colRows.Scan(&field, &colType, &nullable, &key, &defaultValue, &extra)
						columns = append(columns, model.TableColumn{
							ColumnName:    field.String,
							DataType:      colType.String,
							IsNullable:    nullable.String,
							ColumnDefault: defaultValue,
						})
					}
					allSchemas = append(allSchemas, model.TableSchema{
						TableName: table,
						Columns:   columns,
					})
				}

				dbConn.Schemas = allSchemas
			}
		}

		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allDatabaseConnection
}

func (databaseConnectionService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.DatabaseConnectionResponse] {
	paginationResp := model.PaginationResponse[*model.DatabaseConnectionResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allDatabaseConnection []*model.DatabaseConnectionResponse
	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		databaseConnectionEntities, totalItems, err := databaseConnectionService.databaseConnectionRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allDatabaseConnection = mapper.MapDatabaseConnectionEntitiesIntoDatabaseConnectionResponses(databaseConnectionEntities)
		paginationResp = model.PaginationResponse[*model.DatabaseConnectionResponse]{
			Data:        allDatabaseConnection,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: paginationReq.Page,
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

func (databaseConnectionService *ServiceImpl) FindById(ginContext *gin.Context, databaseConnectionId uint64) *model.DatabaseConnectionResponse {
	var databaseConnectionResponse *model.DatabaseConnectionResponse
	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		databaseConnectionEntity, err := databaseConnectionService.databaseConnectionRepository.FindById(gormTransaction, databaseConnectionId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		databaseConnectionResponse = mapper.MapDatabaseConnectionEntityIntoDatabaseConnectionResponse(databaseConnectionEntity)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return databaseConnectionResponse
}

// Create - Membuat databaseConnection baru
func (databaseConnectionService *ServiceImpl) Create(ginContext *gin.Context, createDatabaseConnectionRequest *model.CreateDatabaseConnectionRequest) []*model.DatabaseConnectionResponse {
	valErr := databaseConnectionService.validatorService.ValidateStruct(createDatabaseConnectionRequest)
	databaseConnectionService.validatorService.ParseValidationError(valErr, *createDatabaseConnectionRequest)
	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {

		databaseConnectionEntity := mapper.MapCreateDatabaseConnectionRequestIntoDatabaseConnectionEntity(createDatabaseConnectionRequest)
		err := databaseConnectionService.databaseConnectionRepository.Create(gormTransaction, databaseConnectionEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))

	databaseConnectionResponses := databaseConnectionService.FindAll()
	return databaseConnectionResponses
}

func (databaseConnectionService *ServiceImpl) CreateSchema(ginContext *gin.Context, createDatabaseSchemaRequest *model.CreateDatabaseSchemaRequest) {
	valErr := databaseConnectionService.validatorService.ValidateStruct(createDatabaseSchemaRequest)
	databaseConnectionService.validatorService.ParseValidationError(valErr, *createDatabaseSchemaRequest)
	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		databaseConnection, err := databaseConnectionService.databaseConnectionRepository.FindById(gormTransaction, createDatabaseSchemaRequest.DatabaseConnectionId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		databaseConnectionResponse := mapper.MapDatabaseConnectionEntityIntoDatabaseConnectionResponse(databaseConnection)
		dynamicDatabaseConnection, err := utils.NewDynamicDatabaseConnection(databaseConnectionResponse)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		var stringBuilder strings.Builder

		stringBuilder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", databaseConnectionResponse.Config.Database))

		for i, databaseColumn := range createDatabaseSchemaRequest.Columns {
			stringBuilder.WriteString(fmt.Sprintf("%s %s", databaseColumn.Name, databaseColumn.Type))

			if !databaseColumn.Nullable {
				stringBuilder.WriteString(" NOT NULL")
			}

			if databaseColumn.Primary {
				stringBuilder.WriteString(" PRIMARY KEY")
			}

			if i < len(createDatabaseSchemaRequest.Columns)-1 {
				stringBuilder.WriteString(", ")
			}
		}

		stringBuilder.WriteString(");")

		sqlPayload := stringBuilder.String()
		err = dynamicDatabaseConnection.Exec(sqlPayload).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))

}

func (databaseConnectionService *ServiceImpl) Update(ginContext *gin.Context, updateDatabaseConnectionRequest *model.UpdateDatabaseConnectionRequest) {
	valErr := databaseConnectionService.validatorService.ValidateStruct(updateDatabaseConnectionRequest)
	databaseConnectionService.validatorService.ParseValidationError(valErr, *updateDatabaseConnectionRequest)
	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		databaseConnection, err := databaseConnectionService.databaseConnectionRepository.FindById(gormTransaction, updateDatabaseConnectionRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateDatabaseConnectionRequestIntoDatabaseConnectionEntity(updateDatabaseConnectionRequest, databaseConnection)
		err = databaseConnectionService.databaseConnectionRepository.Update(gormTransaction, databaseConnection)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (databaseConnectionService *ServiceImpl) Delete(ginContext *gin.Context, deleteDatabaseConnectionRequest *model.DeleteDatabaseConnectionRequest) {
	valErr := databaseConnectionService.validatorService.ValidateStruct(deleteDatabaseConnectionRequest)
	databaseConnectionService.validatorService.ParseValidationError(valErr, *deleteDatabaseConnectionRequest)
	err := databaseConnectionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := databaseConnectionService.databaseConnectionRepository.Delete(gormTransaction, deleteDatabaseConnectionRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
