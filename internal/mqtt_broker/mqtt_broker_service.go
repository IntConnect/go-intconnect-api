package mqtt_broker

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createMqttBrokerRequest *model.CreateMqttBrokerRequest)
	FindAll() []*model.MqttBrokerResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.MqttBrokerResponse]
	Update(ginContext *gin.Context, updateMqttBrokerRequest *model.UpdateMqttBrokerRequest)
	Delete(ginContext *gin.Context, deleteMqttBrokerRequest *model.DeleteResourceGeneralRequest)
}
