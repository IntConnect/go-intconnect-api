package database_connection

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	DatabaseConnectionService Service
	viperConfig               *viper.Viper
}

func NewHandler(DatabaseConnectionService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		DatabaseConnectionService: DatabaseConnectionService,
		viperConfig:               viperConfig,
	}
}

func (DatabaseConnectionHandler *Handler) FindAll(ginContext *gin.Context) {
	DatabaseConnectionResponses := DatabaseConnectionHandler.DatabaseConnectionService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Database connection has been fetched", DatabaseConnectionResponses))
}

func (DatabaseConnectionHandler *Handler) FindById(ginContext *gin.Context) {
	DatabaseConnectionId := ginContext.Param("id")
	parsedDatabaseConnectionId, err := strconv.ParseUint(DatabaseConnectionId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, "Supplied param not valid", nil))
	DatabaseConnectionResponses := DatabaseConnectionHandler.DatabaseConnectionService.FindById(ginContext, parsedDatabaseConnectionId)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("DatabaseConnection has been fetched", DatabaseConnectionResponses))
}

func (DatabaseConnectionHandler *Handler) FindAllPagination(ginContext *gin.Context) {
	paginationReq := model.PaginationRequest{
		Page:  1,
		Size:  10,
		Sort:  "id",
		Order: "asc",
	}

	// Bind query parameters to the struct
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	DatabaseConnectionResponses := DatabaseConnectionHandler.DatabaseConnectionService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("DatabaseConnection has been fetched", DatabaseConnectionResponses))
}

func (DatabaseConnectionHandler *Handler) CreateDatabaseConnection(ginContext *gin.Context) {
	var createDatabaseConnectionModel model.CreateDatabaseConnectionRequest
	err := ginContext.ShouldBindBodyWithJSON(&createDatabaseConnectionModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	DatabaseConnectionResponses := DatabaseConnectionHandler.DatabaseConnectionService.Create(ginContext, &createDatabaseConnectionModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("DatabaseConnection has been created", DatabaseConnectionResponses))
}

func (DatabaseConnectionHandler *Handler) CreateDatabaseSchema(ginContext *gin.Context) {
	var createDatabaseSchemaRequest model.CreateDatabaseSchemaRequest
	err := ginContext.ShouldBindBodyWithJSON(&createDatabaseSchemaRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	DatabaseConnectionHandler.DatabaseConnectionService.CreateSchema(ginContext, &createDatabaseSchemaRequest)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("DatabaseConnection has been created", nil))
}

func (DatabaseConnectionHandler *Handler) UpdateDatabaseConnection(ginContext *gin.Context) {
	var updateDatabaseConnectionModel model.UpdateDatabaseConnectionRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateDatabaseConnectionModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	DatabaseConnectionHandler.DatabaseConnectionService.Update(ginContext, &updateDatabaseConnectionModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("DatabaseConnection has been created", nil))
}

func (DatabaseConnectionHandler *Handler) DeleteDatabaseConnection(ginContext *gin.Context) {
	var deleteBomModel model.DeleteDatabaseConnectionRequest
	currencyId := ginContext.Param("id")
	parsedBomId, err := strconv.ParseUint(currencyId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	deleteBomModel.Id = parsedBomId
	DatabaseConnectionHandler.DatabaseConnectionService.Delete(ginContext, &deleteBomModel)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Bom has been updated", nil))
}
