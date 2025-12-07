package role

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll(ginContext *gin.Context) []*model.RoleResponse
	FindById(ginContext *gin.Context, roleId uint64) *model.RoleResponse
	Create(ginContext *gin.Context, createRoleRequest *model.CreateRoleRequest) []*model.RoleResponse
	Update(ginContext *gin.Context, updateRoleRequest *model.UpdateRoleRequest) []*model.RoleResponse
	Delete(ginContext *gin.Context, deleteResourceGeneralRequest *model.DeleteResourceGeneralRequest) []*model.RoleResponse
}
