package protocol_configuration

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	CreateProtocolConfiguration(ginContext *gin.Context)
	DeleteProtocolConfiguration(ginContext *gin.Context)
	UpdateProtocolConfiguration(ginContext *gin.Context)
}
