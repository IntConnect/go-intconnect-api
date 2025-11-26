package report_document_template

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.ReportDocumentTemplate, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.ReportDocumentTemplate, int64, error)
	FindById(gormTransaction *gorm.DB, reportDocumentTemplateId uint64) (*entity.ReportDocumentTemplate, error)
	Create(gormTransaction *gorm.DB, reportDocumentTemplateEntity *entity.ReportDocumentTemplate) error
	Update(gormTransaction *gorm.DB, reportDocumentTemplateEntity *entity.ReportDocumentTemplate) error
	Delete(gormTransaction *gorm.DB, reportDocumentTemplateId uint64) error
}
