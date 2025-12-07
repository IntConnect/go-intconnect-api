package role

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	roleService Service
	viperConfig *viper.Viper
}

func NewHandler(roleService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		roleService: roleService,
		viperConfig: viperConfig,
	}
}

func (roleHandler *Handler) FindAllRole(ginContext *gin.Context) {
	roleResponses := roleHandler.roleService.FindAll(ginContext)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Role has been fetched", roleResponses))
}

func (roleHandler *Handler) CreateRole(ginContext *gin.Context) {
	var createRoleModel model.CreateRoleRequest
	err := ginContext.ShouldBindBodyWithJSON(&createRoleModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	roleResponses := roleHandler.roleService.Create(ginContext, &createRoleModel)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Role has been created", roleResponses))
}

func (roleHandler *Handler) UpdateRole(ginContext *gin.Context) {
	var updateRoleModel model.UpdateRoleRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateRoleModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponses := roleHandler.roleService.Update(ginContext, &updateRoleModel)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Role has been created", paginatedResponses))
}

func (roleHandler *Handler) DeleteRole(ginContext *gin.Context) {
	var deleteRoleModel model.DeleteResourceGeneralRequest
	err := ginContext.ShouldBindBodyWithJSON(&deleteRoleModel)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	roleId := ginContext.Param("id")
	parsedRoleId, err := strconv.ParseUint(roleId, 10, 32)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	deleteRoleModel.Id = parsedRoleId
	paginatedResponses := roleHandler.roleService.Delete(ginContext, &deleteRoleModel)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Role has been deleted", paginatedResponses))
}
