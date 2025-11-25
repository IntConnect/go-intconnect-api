package parameter

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.ParameterResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ParameterResponse]
	Create(ginContext *gin.Context, createParameterRequest *model.CreateParameterRequest) *model.PaginatedResponse[*model.ParameterResponse]
	Update(ginContext *gin.Context, updateParameterRequest *model.UpdateParameterRequest)
	Delete(ginContext *gin.Context, deleteParameterRequest *model.DeleteResourceGeneralRequest)
}
