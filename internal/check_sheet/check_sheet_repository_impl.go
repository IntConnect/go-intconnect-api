package check_sheet

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (checkSheetImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.CheckSheet, error) {
	var checkSheetEntities []*entity.CheckSheet
	err := gormTransaction.Find(&checkSheetEntities).Error
	return checkSheetEntities, err
}

func (checkSheetImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.CheckSheet, int64, error) {

	var checkSheetEntities []*entity.CheckSheet
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.CheckSheet{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("name ILIKE ? OR code ILIKE ?", searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Preload("ReportedByUser").
		Preload("VerifiedByUser").
		Preload("CheckSheetDocumentTemplate").
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&checkSheetEntities).Error; err != nil {
		return nil, 0, err
	}

	return checkSheetEntities, totalItems, nil
}

func (checkSheetImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, checkSheetId uint64) (*entity.CheckSheet, error) {
	var checkSheetEntity entity.CheckSheet
	err := gormTransaction.Model(&entity.CheckSheet{}).
		Preload("ReportedByUser").
		Preload("VerifiedByUser").
		Preload("CheckSheetDocumentTemplate").
		Preload("CheckSheetCheckPoint").
		Preload("CheckSheetCheckPoint.CheckSheetCheckPointValue").
		Where("id = ?", checkSheetId).First(&checkSheetEntity).Error

	return &checkSheetEntity, err
}

func (checkSheetImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.CheckSheet) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (checkSheetImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, checkSheetEntity *entity.CheckSheet) error {
	return gormTransaction.Model(checkSheetEntity).Save(checkSheetEntity).Error
}

func (checkSheetImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.CheckSheet{}).Where("id = ?", id).Delete(&entity.CheckSheet{}).Error
}
