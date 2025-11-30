package audit_log

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (auditLogRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.AuditLog, error) {
	var auditLogEntities []entity.AuditLog
	err := gormTransaction.Find(&auditLogEntities).Error
	return auditLogEntities, err
}

func (auditLogRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.AuditLog, int64, error) {

	var auditLogEntities []entity.AuditLog
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.AuditLog{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("action ILIKE ? OR feature ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Preload("User").
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&auditLogEntities).Error; err != nil {
		return nil, 0, err
	}

	return auditLogEntities, totalItems, nil
}

func (auditLogRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, auditLogId uint64) (*entity.AuditLog, error) {
	var auditLogEntity entity.AuditLog
	err := gormTransaction.Model(&entity.AuditLog{}).Where("id = ?", auditLogId).Find(&auditLogEntity).Error

	return &auditLogEntity, err
}

func (auditLogRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.AuditLog) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}
