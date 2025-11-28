package audit_log

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.AuditLog, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.AuditLog, int64, error)
	FindById(gormTransaction *gorm.DB, auditLogId uint64) (*entity.AuditLog, error)
	Create(gormTransaction *gorm.DB, auditLogEntity *entity.AuditLog) error
}
