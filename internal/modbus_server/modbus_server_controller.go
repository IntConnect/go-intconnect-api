package modbus_server

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllModbusServer(ginContext *gin.Context)
	FindAllModbusServerPagination(ginContext *gin.Context)
	CreateModbusServer(ginContext *gin.Context)
	DeleteModbusServer(ginContext *gin.Context)
	UpdateModbusServer(ginContext *gin.Context)
}
