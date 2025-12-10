package report_document_template

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (reportDocumentTemplateRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.ReportDocumentTemplate, error) {
	var reportDocumentTemplateEntities []entity.ReportDocumentTemplate
	err := gormTransaction.Find(&reportDocumentTemplateEntities).Error
	return reportDocumentTemplateEntities, err
}

func (reportDocumentTemplateRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.ReportDocumentTemplate, int64, error) {

	var reportDocumentTemplateEntities []entity.ReportDocumentTemplate
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.ReportDocumentTemplate{})

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
		Find(&reportDocumentTemplateEntities).Error; err != nil {
		return nil, 0, err
	}

	return reportDocumentTemplateEntities, totalItems, nil
}

func (reportDocumentTemplateRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, reportDocumentTemplateId uint64) (*entity.ReportDocumentTemplate, error) {
	var reportDocumentTemplateEntity entity.ReportDocumentTemplate
	err := gormTransaction.Model(&entity.ReportDocumentTemplate{}).
		Preload("Parameters").
		Where("id = ?", reportDocumentTemplateId).First(&reportDocumentTemplateEntity).Error

	return &reportDocumentTemplateEntity, err
}

func (reportDocumentTemplateRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.ReportDocumentTemplate) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (reportDocumentTemplateRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, reportDocumentTemplateEntity *entity.ReportDocumentTemplate) error {
	return gormTransaction.Model(reportDocumentTemplateEntity).Save(reportDocumentTemplateEntity).Error
}

func (reportDocumentTemplateRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.ReportDocumentTemplate{}).Where("id = ?", id).Delete(&entity.ReportDocumentTemplate{}).Error
}
