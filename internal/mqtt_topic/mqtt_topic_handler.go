package mqtt_topic

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
	mqttTopicService Service
	viperConfig      *viper.Viper
}

func NewHandler(mqttTopicService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		mqttTopicService: mqttTopicService,
		viperConfig:      viperConfig,
	}
}

func (mqttTopicHandler *Handler) FindAllMqttTopic(ginContext *gin.Context) {
	mqttTopicResponses := mqttTopicHandler.mqttTopicService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("MQTT Topic has been fetched", mqttTopicResponses))
}

func (mqttTopicHandler *Handler) FindAllMqttTopicPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := mqttTopicHandler.mqttTopicService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttTopicHandler *Handler) CreateMqttTopic(ginContext *gin.Context) {
	var createMqttTopicModel model.CreateMqttTopicRequest
	err := ginContext.ShouldBindBodyWithJSON(&createMqttTopicModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := mqttTopicHandler.mqttTopicService.Create(ginContext, &createMqttTopicModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttTopicHandler *Handler) FindDependencyMqttTopic(ginContext *gin.Context) {
	mqttTopicDependency := mqttTopicHandler.mqttTopicService.FindDependency()
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse(model.RESPONSE_SUCCESS, mqttTopicDependency))

}

func (mqttTopicHandler *Handler) UpdateMqttTopic(ginContext *gin.Context) {
	var updateMqttTopicModel model.UpdateMqttTopicRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateMqttTopicModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	mqttTopicId := ginContext.Param("id")
	parsedMqttTopicId, err := strconv.ParseUint(mqttTopicId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	updateMqttTopicModel.Id = parsedMqttTopicId
	paginatedResponse := mqttTopicHandler.mqttTopicService.Update(ginContext, &updateMqttTopicModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (mqttTopicHandler *Handler) DeleteMqttTopic(ginContext *gin.Context) {
	var deleteMqttTopicModel model.DeleteResourceGeneralRequest
	mqttTopicId := ginContext.Param("id")
	parsedMqttTopicId, err := strconv.ParseUint(mqttTopicId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteMqttTopicModel.Id = parsedMqttTopicId
	paginatedResponse := mqttTopicHandler.mqttTopicService.Delete(ginContext, &deleteMqttTopicModel)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
