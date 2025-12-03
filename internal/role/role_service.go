package role

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ginContext *gin.Context, createRoleRequest *model.CreateRoleRequest)
	FindAll(ginContext *gin.Context) []*model.RoleResponse
	FindById(ginContext *gin.Context, roleId uint64) *model.RoleResponse
	Update(ginContext *gin.Context, updateRoleRequest *model.UpdateRoleRequest)
	Delete(ginContext *gin.Context, deleteResourceGeneralRequest *model.DeleteResourceGeneralRequest)
}
