package check_sheet_document_template_parameter

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.CheckSheetDocumentTemplateParameter, error) {
	var checkSheetDocumentTemplateParameterEntities []entity.CheckSheetDocumentTemplateParameter
	err := gormTransaction.Find(&checkSheetDocumentTemplateParameterEntities).Error
	return checkSheetDocumentTemplateParameterEntities, err
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.CheckSheetDocumentTemplateParameter, int64, error) {

	var checkSheetDocumentTemplateParameterEntities []entity.CheckSheetDocumentTemplateParameter
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.CheckSheetDocumentTemplateParameter{})

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
		Find(&checkSheetDocumentTemplateParameterEntities).Error; err != nil {
		return nil, 0, err
	}

	return checkSheetDocumentTemplateParameterEntities, totalItems, nil
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterId uint64) (*entity.CheckSheetDocumentTemplateParameter, error) {
	var checkSheetDocumentTemplateParameterEntity entity.CheckSheetDocumentTemplateParameter
	err := gormTransaction.Model(&entity.CheckSheetDocumentTemplateParameter{}).
		Where("id = ?", checkSheetDocumentTemplateParameterId).
		Find(&checkSheetDocumentTemplateParameterEntity).Error

	return &checkSheetDocumentTemplateParameterEntity, err
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterEntities []*entity.CheckSheetDocumentTemplateParameter) error {
	return gormTransaction.Model(checkSheetDocumentTemplateParameterEntities).Create(checkSheetDocumentTemplateParameterEntities).Error
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) DeleteBatch(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterEntities []*entity.CheckSheetDocumentTemplateParameter) error {
	return gormTransaction.Delete(checkSheetDocumentTemplateParameterEntities).Error
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) DeleteBatchById(gormTransaction *gorm.DB, checkSheetDocumentTemplateParameterId []uint64) error {
	return gormTransaction.Model(entity.CheckSheetDocumentTemplateParameter{}).Where("check_sheet_id IN ?", checkSheetDocumentTemplateParameterId).Delete(entity.CheckSheetDocumentTemplateParameter{}).Error
}

func (checkSheetDocumentTemplateParameterRepositoryImpl *RepositoryImpl) DeleteBatchByCheckSheetId(gormTransaction *gorm.DB, checkSheetId uint64) error {
	return gormTransaction.Model(entity.CheckSheetDocumentTemplateParameter{}).Where("check_sheet_id = ?", checkSheetId).Delete(entity.CheckSheetDocumentTemplateParameter{}).Error
}
