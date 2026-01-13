package check_sheet_check_point

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.CheckSheetCheckPoint, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.CheckSheetCheckPoint, int64, error)
	FindById(gormTransaction *gorm.DB, checkSheetCheckPointId uint64) (*entity.CheckSheetCheckPoint, error)
	Create(gormTransaction *gorm.DB, checkSheetCheckPointEntity *entity.CheckSheetCheckPoint) error
	CreateBatch(gormTransaction *gorm.DB, checkSheetCheckPointEntities []*entity.CheckSheetCheckPoint) error
	Delete(gormTransaction *gorm.DB, checkSheetCheckPointId uint64) error
	DeleteBatchById(gormTransaction *gorm.DB, checkSheetId uint64) error
}
