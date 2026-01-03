package register

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createRegisterRequest *model.CreateRegisterRequest) *model.PaginatedResponse[*model.RegisterResponse]
	FindAll() []*model.RegisterResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.RegisterResponse]
	Update(ginContext *gin.Context, updateRegisterRequest *model.UpdateRegisterRequest) *model.PaginatedResponse[*model.RegisterResponse]
	Delete(ginContext *gin.Context, deleteRegisterRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.RegisterResponse]
	FindDependency() *model.RegisterDependency
}
