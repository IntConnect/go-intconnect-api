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
		beforeEntity interface{},
		afterEntity interface{},
		description string,
	) error
}
