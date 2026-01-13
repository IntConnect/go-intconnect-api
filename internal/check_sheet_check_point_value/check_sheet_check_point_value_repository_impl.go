package check_sheet_check_point

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (checkSheetCheckPointValueRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, checkSheetCheckPointValueEntity *entity.CheckSheetCheckPointValue) error {
	return gormTransaction.Model(checkSheetCheckPointValueEntity).Create(checkSheetCheckPointValueEntity).Error
}

func (checkSheetCheckPointValueRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, checkSheetCheckPointValueEntities []*entity.CheckSheetCheckPointValue) error {
	return gormTransaction.Model(checkSheetCheckPointValueEntities).Create(checkSheetCheckPointValueEntities).Error
}

func (checkSheetCheckPointValueRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.CheckSheetCheckPointValue{}).Where("id = ?", id).Delete(entity.CheckSheetCheckPointValue{}).Error
}

func (checkSheetCheckPointValueRepositoryImpl *RepositoryImpl) DeleteBatchById(gormTransaction *gorm.DB, checkSheetId uint64) error {
	return gormTransaction.Model(entity.CheckSheetCheckPointValue{}).Where("check_sheet_id = ?", checkSheetId).Delete(entity.CheckSheetCheckPointValue{}).Error
}
