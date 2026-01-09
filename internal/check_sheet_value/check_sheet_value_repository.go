package check_sheet_value

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.CheckSheetValue, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.CheckSheetValue, int64, error)
	FindById(gormTransaction *gorm.DB, checkSheetValueId uint64) (*entity.CheckSheetValue, error)
	Create(gormTransaction *gorm.DB, checkSheetValueEntity *entity.CheckSheetValue) error
	CreateBatch(gormTransaction *gorm.DB, checkSheetValueEntities []*entity.CheckSheetValue) error
	Delete(gormTransaction *gorm.DB, checkSheetValueId uint64) error
	DeleteBatchById(gormTransaction *gorm.DB, checkSheetId uint64) error
}
