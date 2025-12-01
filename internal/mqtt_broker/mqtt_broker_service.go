package mqtt_broker

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createMqttBrokerRequest *model.CreateMqttBrokerRequest) *model.PaginatedResponse[*model.MqttBrokerResponse]
	FindAll() []*model.MqttBrokerResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MqttBrokerResponse]
	Update(ginContext *gin.Context, updateMqttBrokerRequest *model.UpdateMqttBrokerRequest) *model.PaginatedResponse[*model.MqttBrokerResponse]
	Delete(ginContext *gin.Context, deleteMqttBrokerRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.MqttBrokerResponse]
}
