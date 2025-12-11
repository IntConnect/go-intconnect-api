package breakdown

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.BreakdownResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.BreakdownResponse]
	Create(ginContext *gin.Context, createBreakdownRequest *model.CreateBreakdownRequest) *model.PaginatedResponse[*model.BreakdownResponse]
	Update(ginContext *gin.Context, updateBreakdownRequest *model.UpdateBreakdownRequest) *model.PaginatedResponse[*model.BreakdownResponse]
	Delete(ginContext *gin.Context, deleteBreakdownRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.BreakdownResponse]
}
