package facility

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	CreateFacility(ginContext *gin.Context)
	DeleteFacility(ginContext *gin.Context)
	UpdateFacility(ginContext *gin.Context)
}
