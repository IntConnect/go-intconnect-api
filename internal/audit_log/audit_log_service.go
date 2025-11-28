package audit_log

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.AuditLogResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.AuditLogResponse]
	Create(ginContext *gin.Context, createAuditLogRequest *model.CreateAuditLogRequest) *model.PaginatedResponse[*model.AuditLogResponse]
}
