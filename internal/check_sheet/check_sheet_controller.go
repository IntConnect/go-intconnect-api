package check_sheet

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllCheckSheet(ginContext *gin.Context)
	FindAllCheckSheetPagination(ginContext *gin.Context)
	FindCheckSheetById(ginContext *gin.Context)
	CreateCheckSheet(ginContext *gin.Context)
	UpdateCheckSheet(ginContext *gin.Context)
	DeleteCheckSheet(ginContext *gin.Context)
}
