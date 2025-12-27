package check_sheet_document_template

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (checkSheetDocumentTemplateRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.CheckSheetDocumentTemplate, error) {
	var checkSheetDocumentTemplateEntities []*entity.CheckSheetDocumentTemplate
	err := gormTransaction.
		Preload("Parameters").
		Find(&checkSheetDocumentTemplateEntities).Error
	return checkSheetDocumentTemplateEntities, err
}

func (checkSheetDocumentTemplateRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.CheckSheetDocumentTemplate, int64, error) {

	var checkSheetDocumentTemplateEntities []*entity.CheckSheetDocumentTemplate
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.CheckSheetDocumentTemplate{})

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
		Preload("Parameters").
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&checkSheetDocumentTemplateEntities).Error; err != nil {
		return nil, 0, err
	}

	return checkSheetDocumentTemplateEntities, totalItems, nil
}

func (checkSheetDocumentTemplateRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, checkSheetDocumentTemplateId uint64) (*entity.CheckSheetDocumentTemplate, error) {
	var checkSheetDocumentTemplateEntity entity.CheckSheetDocumentTemplate
	err := gormTransaction.Model(&entity.CheckSheetDocumentTemplate{}).
		Preload("Parameters").
		Where("id = ?", checkSheetDocumentTemplateId).First(&checkSheetDocumentTemplateEntity).Error

	return &checkSheetDocumentTemplateEntity, err
}

func (checkSheetDocumentTemplateRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.CheckSheetDocumentTemplate) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (checkSheetDocumentTemplateRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, checkSheetDocumentTemplateEntity *entity.CheckSheetDocumentTemplate) error {
	return gormTransaction.Model(checkSheetDocumentTemplateEntity).Save(checkSheetDocumentTemplateEntity).Error
}

func (checkSheetDocumentTemplateRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.CheckSheetDocumentTemplate{}).Where("id = ?", id).Delete(&entity.CheckSheetDocumentTemplate{}).Error
}
