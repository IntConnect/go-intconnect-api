package breakdown

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllBreakdown(ginContext *gin.Context)
	FindAllBreakdownPagination(ginContext *gin.Context)
	CreateBreakdown(ginContext *gin.Context)
	DeleteBreakdown(ginContext *gin.Context)
	UpdateBreakdown(ginContext *gin.Context)
}
