package check_sheet_document_template

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllCheckSheetDocumentTemplate(ginContext *gin.Context)
	FindAllCheckSheetDocumentTemplatePagination(ginContext *gin.Context)
	CreateCheckSheetDocumentTemplate(ginContext *gin.Context)
	UpdateCheckSheetDocumentTemplate(ginContext *gin.Context)
	DeleteCheckSheetDocumentTemplate(ginContext *gin.Context)
}
