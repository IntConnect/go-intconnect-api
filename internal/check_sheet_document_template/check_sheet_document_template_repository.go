package check_sheet_document_template

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.CheckSheetDocumentTemplate, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.CheckSheetDocumentTemplate, int64, error)
	FindById(gormTransaction *gorm.DB, checkSheetDocumentTemplateId uint64) (*entity.CheckSheetDocumentTemplate, error)
	Create(gormTransaction *gorm.DB, checkSheetDocumentTemplateEntity *entity.CheckSheetDocumentTemplate) error
	Update(gormTransaction *gorm.DB, checkSheetDocumentTemplateEntity *entity.CheckSheetDocumentTemplate) error
	Delete(gormTransaction *gorm.DB, checkSheetDocumentTemplateId uint64) error
}
