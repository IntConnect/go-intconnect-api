package telemetry

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.Telemetry, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Telemetry, int64, error)
	FindById(gormTransaction *gorm.DB, telemetryId uint64) (*entity.Telemetry, error)
	Create(gormTransaction *gorm.DB, telemetryEntity *entity.Telemetry) error
	CreateBatch(gormTransaction *gorm.DB, telemetryEntities []*entity.Telemetry) error
	Update(gormTransaction *gorm.DB, telemetryEntity *entity.Telemetry) error
	Delete(gormTransaction *gorm.DB, telemetryId uint64) error
}
