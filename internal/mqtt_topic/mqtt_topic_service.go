package mqtt_topic

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.MqttTopicResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MqttTopicResponse]
	Create(ginContext *gin.Context, createMqttTopicRequest *model.CreateMqttTopicRequest) *model.PaginatedResponse[*model.MqttTopicResponse]
	Update(ginContext *gin.Context, updateMqttTopicRequest *model.UpdateMqttTopicRequest)
	Delete(ginContext *gin.Context, deleteMqttTopicRequest *model.DeleteResourceGeneralRequest)
}
