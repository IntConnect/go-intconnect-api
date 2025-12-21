package telemetry

import (
	"go-intconnect-api/internal/entity"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAllFilter(gormTransaction *gorm.DB, searchedParameterIds []uint64, intervalVal string, startDate, endDate time.Time) ([]*entity.TelemetryQuery, error)
	CreateBatch(gormTransaction *gorm.DB, telemetryEntities []*entity.Telemetry) error
}
