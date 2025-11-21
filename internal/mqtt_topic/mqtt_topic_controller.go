package mqtt_topic

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllMqttTopic(ginContext *gin.Context)
	FindAllMqttTopicPagination(ginContext *gin.Context)
	CreateMqttTopic(ginContext *gin.Context)
	DeleteMqttTopic(ginContext *gin.Context)
	UpdateMqttTopic(ginContext *gin.Context)
}
