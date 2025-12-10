package report_document_template

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllReportDocumentTemplate(ginContext *gin.Context)
	FindAllReportDocumentTemplatePagination(ginContext *gin.Context)
	CreateReportDocumentTemplate(ginContext *gin.Context)
	UpdateReportDocumentTemplate(ginContext *gin.Context)
	DeleteReportDocumentTemplate(ginContext *gin.Context)
}
