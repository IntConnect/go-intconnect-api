package parameter

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllParameter(ginContext *gin.Context)
	FindDependencyParameter(context *gin.Context)
	FindAllParameterPagination(ginContext *gin.Context)
	CreateParameter(ginContext *gin.Context)
	UpdateParameter(ginContext *gin.Context)
	UpdateParameterOperation(ginContext *gin.Context)
	DeleteParameter(ginContext *gin.Context)
}
