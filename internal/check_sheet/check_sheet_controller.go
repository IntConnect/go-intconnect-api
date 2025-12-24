package check_sheet_document_template

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllCheckSheet(ginContext *gin.Context)
	FindAllCheckSheetPagination(ginContext *gin.Context)
	CreateCheckSheet(ginContext *gin.Context)
	UpdateCheckSheet(ginContext *gin.Context)
	DeleteCheckSheet(ginContext *gin.Context)
}
