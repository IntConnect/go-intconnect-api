package pipeline

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	CreatePipeline(ginContext *gin.Context)
	DeletePipeline(ginContext *gin.Context)
	UpdatePipeline(ginContext *gin.Context)
}
