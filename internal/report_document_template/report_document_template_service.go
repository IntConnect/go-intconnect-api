package report_document_template

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.ReportDocumentTemplateResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse]
	Create(ginContext *gin.Context, createReportDocumentTemplateRequest *model.CreateReportDocumentTemplateRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse]
	Update(ginContext *gin.Context, updateReportDocumentTemplateRequest *model.UpdateReportDocumentTemplateRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse]
	Delete(ginContext *gin.Context, deleteReportDocumentTemplateRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.ReportDocumentTemplateResponse]
}
