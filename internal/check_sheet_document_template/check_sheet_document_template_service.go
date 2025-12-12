package check_sheet_document_template

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.CheckSheetDocumentTemplateResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse]
	Create(ginContext *gin.Context, createCheckSheetDocumentTemplateRequest *model.CreateCheckSheetDocumentTemplateRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse]
	Update(ginContext *gin.Context, updateCheckSheetDocumentTemplateRequest *model.UpdateCheckSheetDocumentTemplateRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse]
	Delete(ginContext *gin.Context, deleteCheckSheetDocumentTemplateRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.CheckSheetDocumentTemplateResponse]
}
