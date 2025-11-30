package audit_log

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	auditLogService Service
	viperConfig     *viper.Viper
}

func NewHandler(auditLogService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		auditLogService: auditLogService,
		viperConfig:     viperConfig,
	}
}

func (auditLogHandler *Handler) FindAllAuditLog(ginContext *gin.Context) {
	auditLogResponses := auditLogHandler.auditLogService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("AuditLog has been fetched", auditLogResponses))
}

func (auditLogHandler *Handler) FindAllAuditLogPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := auditLogHandler.auditLogService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}
