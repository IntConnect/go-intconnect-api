package database_connection

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createDatabaseConnectionRequest *model.CreateDatabaseConnectionRequest) []*model.DatabaseConnectionResponse
	FindAll() []*model.DatabaseConnectionResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.DatabaseConnectionResponse]
	Update(ginContext *gin.Context, updateDatabaseConnectionRequest *model.UpdateDatabaseConnectionRequest)
	Delete(ginContext *gin.Context, deleteDatabaseConnectionRequest *model.DeleteDatabaseConnectionRequest)
	FindById(ginContext *gin.Context, protocolConfigurationId uint64) *model.DatabaseConnectionResponse
	CreateSchema(ginContext *gin.Context, createDatabaseSchemaRequest *model.CreateDatabaseSchemaRequest)
}
