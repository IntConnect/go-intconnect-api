package mqtt_topic

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.MqttTopicResponse
	FindDependency() *model.MqttTopicDependency
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MqttTopicResponse]
	Create(ginContext *gin.Context, createMqttTopicRequest *model.CreateMqttTopicRequest) *model.PaginatedResponse[*model.MqttTopicResponse]
	Update(ginContext *gin.Context, updateMqttTopicRequest *model.UpdateMqttTopicRequest) *model.PaginatedResponse[*model.MqttTopicResponse]
	Delete(ginContext *gin.Context, deleteMqttTopicRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.MqttTopicResponse]
}
