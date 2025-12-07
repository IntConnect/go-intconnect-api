package modbus_server

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
	modbusServerService Service
	viperConfig         *viper.Viper
}

func NewHandler(modbusServerService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		modbusServerService: modbusServerService,
		viperConfig:         viperConfig,
	}
}

func (modbusServerHandler *Handler) FindAllModbusServer(ginContext *gin.Context) {
	modbusServerResponses := modbusServerHandler.modbusServerService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("MQTT Topic has been fetched", modbusServerResponses))
}

func (modbusServerHandler *Handler) FindAllModbusServerPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := modbusServerHandler.modbusServerService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (modbusServerHandler *Handler) CreateModbusServer(ginContext *gin.Context) {
	var createModbusServerModel model.CreateModbusServerRequest
	err := ginContext.ShouldBindBodyWithJSON(&createModbusServerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := modbusServerHandler.modbusServerService.Create(ginContext, &createModbusServerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (modbusServerHandler *Handler) UpdateModbusServer(ginContext *gin.Context) {
	var updateModbusServerModel model.UpdateModbusServerRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateModbusServerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	modbusServerId := ginContext.Param("id")
	parsedModbusServerId, err := strconv.ParseUint(modbusServerId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateModbusServerModel.Id = parsedModbusServerId
	paginatedResponse := modbusServerHandler.modbusServerService.Update(ginContext, &updateModbusServerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (modbusServerHandler *Handler) DeleteModbusServer(ginContext *gin.Context) {
	var deleteModbusServerModel model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteModbusServerModel)
	modbusServerId := ginContext.Param("id")
	parsedModbusServerId, err := strconv.ParseUint(modbusServerId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteModbusServerModel.Id = parsedModbusServerId
	paginatedResponse := modbusServerHandler.modbusServerService.Delete(ginContext, &deleteModbusServerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
