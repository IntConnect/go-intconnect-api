package parameter

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll(parameterFilterRequest *model.ParameterFilterRequest) []*model.ParameterResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ParameterResponse]
	FindDependencyParameter() *model.ParameterDependency
	FindById(ginContext *gin.Context, parameterId uint64) *model.ParameterResponse
	Create(ginContext *gin.Context, createParameterRequest *model.CreateParameterRequest) *model.PaginatedResponse[*model.ParameterResponse]
	Update(ginContext *gin.Context, updateParameterRequest *model.UpdateParameterRequest)
	UpdateOperation(ginContext *gin.Context, updateParameterRequest *model.ManageParameterOperationRequest) *model.PaginatedResponse[*model.ParameterResponse]
	Delete(ginContext *gin.Context, deleteParameterRequest *model.DeleteResourceGeneralRequest)
}
