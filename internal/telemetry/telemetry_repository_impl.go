package telemetry

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (telemetryRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Telemetry, error) {
	var telemetryEntities []entity.Telemetry
	err := gormTransaction.Find(&telemetryEntities).Error
	return telemetryEntities, err
}

func (telemetryRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.Telemetry, int64, error) {

	var telemetryEntities []entity.Telemetry
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.Telemetry{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("telemetryname ILIKE ? OR email ILIKE ? OR name ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Preload("TelemetryGroup", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id,name")
		}).
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&telemetryEntities).Error; err != nil {
		return nil, 0, err
	}

	return telemetryEntities, totalItems, nil
}

func (telemetryRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, telemetryId uint64) (*entity.Telemetry, error) {
	var telemetryEntity entity.Telemetry
	err := gormTransaction.Model(&entity.Telemetry{}).
		Preload("TelemetryGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", telemetryId).Find(&telemetryEntity).Error

	return &telemetryEntity, err
}

func (telemetryRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, telemetryEntity *entity.Telemetry) error {
	return gormTransaction.Model(telemetryEntity).Create(telemetryEntity).Error
}
func (telemetryRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, telemetryEntities []*entity.Telemetry) error {
	return gormTransaction.Model(telemetryEntities).Create(telemetryEntities).Error
}

func (telemetryRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, telemetryEntity *entity.Telemetry) error {
	return gormTransaction.Model(telemetryEntity).Save(telemetryEntity).Error
}

func (telemetryRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Telemetry{}).Where("id = ?", id).Delete(entity.Telemetry{}).Error
}
