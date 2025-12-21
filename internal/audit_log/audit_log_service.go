package audit_log

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.AuditLogResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.AuditLogResponse]
	Create(ginContext *gin.Context, createAuditLogRequest *model.CreateAuditLogRequest)

	Record(
		ginContext *gin.Context,
		actionType string,
		featureType string,
		auditLogPayload model.AuditLogPayload,
	) error
	Build(
		beforeEntity interface{},
		afterEntity interface{},
		relation map[string]map[string][]uint64,
		description string,
	) model.AuditLogPayload
}
