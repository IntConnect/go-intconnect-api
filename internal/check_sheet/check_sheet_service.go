package check_sheet

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.CheckSheetResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.CheckSheetResponse]
	FindById(ginContext *gin.Context, checkSheetId uint64) *model.CheckSheetResponse
	Create(ginContext *gin.Context, createCheckSheetRequest *model.CreateCheckSheetRequest) *model.PaginatedResponse[*model.CheckSheetResponse]
	Approval(ginContext *gin.Context, approvalCheckSheetRequest *model.ApprovalCheckSheet) *model.PaginatedResponse[*model.CheckSheetResponse]
	Update(ginContext *gin.Context, updateCheckSheetRequest *model.UpdateCheckSheetRequest) *model.PaginatedResponse[*model.CheckSheetResponse]
	Delete(ginContext *gin.Context, deleteCheckSheetRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.CheckSheetResponse]
}
