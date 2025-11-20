package mqtt_broker

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllMqttBroker(ginContext *gin.Context)
	FindAllMqttBrokerPagination(ginContext *gin.Context)
	CreateMqttBroker(ginContext *gin.Context)
	DeleteMqttBroker(ginContext *gin.Context)
	UpdateMqttBroker(ginContext *gin.Context)
}
