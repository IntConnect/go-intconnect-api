package check_sheet_document_template_parameter

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.CheckSheetDocumentTemplateParameter, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.CheckSheetDocumentTemplateParameter, int64, error)
	FindById(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterId uint64) (*entity.CheckSheetDocumentTemplateParameter, error)
	CreateBatch(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterEntities []*entity.CheckSheetDocumentTemplateParameter) error
	DeleteBatch(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterId []*entity.CheckSheetDocumentTemplateParameter) error
	DeleteBatchById(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterId []uint64) error
	DeleteBatchByCheckSheetId(gormTransaction *gorm.DB, checkSheetId uint64) error
}
