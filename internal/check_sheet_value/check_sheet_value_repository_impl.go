package check_sheet_value

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.CheckSheetValue, error) {
	var checkSheetValueEntities []entity.CheckSheetValue
	err := gormTransaction.Find(&checkSheetValueEntities).Error
	return checkSheetValueEntities, err
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.CheckSheetValue, int64, error) {

	var checkSheetValueEntities []entity.CheckSheetValue
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.CheckSheetValue{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("name ILIKE ? OR email ILIKE ? OR name ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&checkSheetValueEntities).Error; err != nil {
		return nil, 0, err
	}

	return checkSheetValueEntities, totalItems, nil
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, checkSheetValueId uint64) (*entity.CheckSheetValue, error) {
	var checkSheetValueEntity entity.CheckSheetValue
	err := gormTransaction.Model(&entity.CheckSheetValue{}).
		Where("id = ?", checkSheetValueId).
		Find(&checkSheetValueEntity).Error

	return &checkSheetValueEntity, err
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, checkSheetValueEntity *entity.CheckSheetValue) error {
	return gormTransaction.Model(checkSheetValueEntity).Create(checkSheetValueEntity).Error
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, checkSheetValueEntities []*entity.CheckSheetValue) error {
	return gormTransaction.Model(checkSheetValueEntities).Create(checkSheetValueEntities).Error
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.CheckSheetValue{}).Where("id = ?", id).Delete(entity.CheckSheetValue{}).Error
}

func (checkSheetValueRepositoryImpl *RepositoryImpl) DeleteBatchById(gormTransaction *gorm.DB, checkSheetId uint64) error {
	return gormTransaction.Model(entity.CheckSheetValue{}).Where("check_sheet_id = ?", checkSheetId).Delete(entity.CheckSheetValue{}).Error
}
