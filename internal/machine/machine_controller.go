package machine

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllMachine(ginContext *gin.Context)
	FindAllMachinePagination(ginContext *gin.Context)
	CreateMachine(ginContext *gin.Context)
	DeleteMachine(ginContext *gin.Context)
	UpdateMachine(ginContext *gin.Context)
}
