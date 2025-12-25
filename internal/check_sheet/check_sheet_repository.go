package check_sheet

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.CheckSheet, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.CheckSheet, int64, error)
	FindById(gormTransaction *gorm.DB, checkSheetId uint64) (*entity.CheckSheet, error)
	Create(gormTransaction *gorm.DB, checkSheetEntity *entity.CheckSheet) error
	Update(gormTransaction *gorm.DB, checkSheetEntity *entity.CheckSheet) error
	Delete(gormTransaction *gorm.DB, checkSheetId uint64) error
}
