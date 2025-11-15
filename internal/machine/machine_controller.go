package machine

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	CreateMachine(ginContext *gin.Context)
	DeleteMachine(ginContext *gin.Context)
	UpdateMachine(ginContext *gin.Context)
}
