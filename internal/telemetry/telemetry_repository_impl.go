package telemetry

import (
	"go-intconnect-api/internal/entity"
	"time"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (telemetryRepositoryImpl *RepositoryImpl) FindAllFilter(gormTransaction *gorm.DB, searchedParameterIds []uint64, intervalVal string, startDate, endDate time.Time) ([]*entity.TelemetryQuery, error) {
	sqlQuery := `
SELECT
	id,
  bucket,
  parameter_id,
  last_value
FROM (
    SELECT 
        id,
      time_bucket_gapfill(?::interval, timestamp) AS bucket,
      parameter_id,
      last(value, timestamp) AS last_value
    FROM telemetries
    WHERE parameter_id IN (?)
      AND timestamp BETWEEN ? AND ?
    GROUP BY id, bucket, parameter_id
) q
ORDER BY bucket;
`

	var telemetryEntities []*entity.TelemetryQuery
	err := gormTransaction.Raw(sqlQuery, intervalVal, searchedParameterIds, startDate, endDate).Scan(&telemetryEntities).Error
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

func (telemetryRepositoryImpl *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, telemetryEntities []*entity.Telemetry) error {
	return gormTransaction.Model(telemetryEntities).Create(telemetryEntities).Error
}
