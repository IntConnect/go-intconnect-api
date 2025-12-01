package mqtt_broker

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
	mqttBrokerService Service
	viperConfig       *viper.Viper
}

func NewHandler(mqttBrokerService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		mqttBrokerService: mqttBrokerService,
		viperConfig:       viperConfig,
	}
}

func (mqttBrokerHandler *Handler) FindAllMqttBroker(ginContext *gin.Context) {
	mqttBrokerResponses := mqttBrokerHandler.mqttBrokerService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("MqttBroker has been fetched", mqttBrokerResponses))
}

func (mqttBrokerHandler *Handler) FindAllMqttBrokerPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttBrokerHandler *Handler) CreateMqttBroker(ginContext *gin.Context) {
	var createMqttBrokerModel model.CreateMqttBrokerRequest
	err := ginContext.ShouldBindBodyWithJSON(&createMqttBrokerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.Create(ginContext, &createMqttBrokerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttBrokerHandler *Handler) UpdateMqttBroker(ginContext *gin.Context) {
	var updateMqttBrokerModel model.UpdateMqttBrokerRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateMqttBrokerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	mqttBrokerId := ginContext.Param("id")
	parsedMqttBrokerId, err := strconv.ParseUint(mqttBrokerId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateMqttBrokerModel.Id = parsedMqttBrokerId
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.Update(ginContext, &updateMqttBrokerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttBrokerHandler *Handler) DeleteMqttBroker(ginContext *gin.Context) {
	var deleteMqttBrokerModel model.DeleteResourceGeneralRequest
	mqttBrokerId := ginContext.Param("id")
	parsedMqttBrokerId, err := strconv.ParseUint(mqttBrokerId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))

	err = ginContext.ShouldBindBodyWithJSON(&deleteMqttBrokerModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteMqttBrokerModel.Id = parsedMqttBrokerId
	paginatedResponse := mqttBrokerHandler.mqttBrokerService.Delete(ginContext, &deleteMqttBrokerModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
