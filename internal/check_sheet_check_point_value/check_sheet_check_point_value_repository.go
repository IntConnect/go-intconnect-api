package check_sheet_check_point

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Create(gormTransaction *gorm.DB, checkSheetCheckPointValueEntity *entity.CheckSheetCheckPointValue) error
	CreateBatch(gormTransaction *gorm.DB, checkSheetCheckPointValueEntities []*entity.CheckSheetCheckPointValue) error
	Delete(gormTransaction *gorm.DB, checkSheetCheckPointValueId uint64) error
	DeleteBatchById(gormTransaction *gorm.DB, checkSheetId uint64) error
}
