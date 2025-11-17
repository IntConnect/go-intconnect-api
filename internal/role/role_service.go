package role

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createRoleRequest *model.CreateRoleRequest)
	FindAll() []*model.RoleResponse
	FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.RoleResponse]
	Update(ginContext *gin.Context, updateRoleRequest *model.UpdateRoleRequest)
	Delete(ginContext *gin.Context, deleteResourceGeneralRequest *model.DeleteResourceGeneralRequest)
}
