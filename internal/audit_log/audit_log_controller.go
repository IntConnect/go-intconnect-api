package audit_log

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllAuditLog(ginContext *gin.Context)
	FindAllAuditLogPagination(ginContext *gin.Context)
}
