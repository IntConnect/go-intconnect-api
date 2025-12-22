package facility

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllFacility(ginContext *gin.Context)
	FindAllFacilityPagination(ginContext *gin.Context)
	FindMinimalAllFacility(context *gin.Context)
	FindFacilityById(ginContext *gin.Context)
	CreateFacility(ginContext *gin.Context)
	DeleteFacility(ginContext *gin.Context)
	UpdateFacility(ginContext *gin.Context)
}
