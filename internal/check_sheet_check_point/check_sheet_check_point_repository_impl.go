package check_sheet_check_point

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.CheckSheetCheckPoint, error) {
	var checkSheetCheckPointEntities []entity.CheckSheetCheckPoint
	err := gormTransaction.Find(&checkSheetCheckPointEntities).Error
	return checkSheetCheckPointEntities, err
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.CheckSheetCheckPoint, int64, error) {

	var checkSheetCheckPointEntities []entity.CheckSheetCheckPoint
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.CheckSheetCheckPoint{})

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
		Find(&checkSheetCheckPointEntities).Error; err != nil {
		return nil, 0, err
	}

	return checkSheetCheckPointEntities, totalItems, nil
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, checkSheetCheckPointId uint64) (*entity.CheckSheetCheckPoint, error) {
	var checkSheetCheckPointEntity entity.CheckSheetCheckPoint
	err := gormTransaction.Model(&entity.CheckSheetCheckPoint{}).
		Where("id = ?", checkSheetCheckPointId).
		Find(&checkSheetCheckPointEntity).Error

	return &checkSheetCheckPointEntity, err
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, checkSheetCheckPointEntity *entity.CheckSheetCheckPoint) error {
	return gormTransaction.Model(checkSheetCheckPointEntity).Create(checkSheetCheckPointEntity).Error
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, checkSheetCheckPointEntities []*entity.CheckSheetCheckPoint) error {
	return gormTransaction.Model(checkSheetCheckPointEntities).Create(checkSheetCheckPointEntities).Error
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.CheckSheetCheckPoint{}).Where("id = ?", id).Delete(entity.CheckSheetCheckPoint{}).Error
}

func (checkSheetCheckPointRepositoryImpl *RepositoryImpl) DeleteBatchById(gormTransaction *gorm.DB, checkSheetId uint64) error {
	return gormTransaction.Model(entity.CheckSheetCheckPoint{}).Where("check_sheet_id = ?", checkSheetId).Delete(entity.CheckSheetCheckPoint{}).Error
}
