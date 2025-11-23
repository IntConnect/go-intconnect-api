package permission

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	permissionService Service
	viperConfig       *viper.Viper
}

func NewHandler(permissionService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		permissionService: permissionService,
		viperConfig:       viperConfig,
	}
}

func (permissionHandler *Handler) FindAllPermission(ginContext *gin.Context) {
	permissionResponses := permissionHandler.permissionService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Permission has been fetched", permissionResponses))
}

func (permissionHandler *Handler) FindAllPermissionPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	paginatedResponse := permissionHandler.permissionService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
