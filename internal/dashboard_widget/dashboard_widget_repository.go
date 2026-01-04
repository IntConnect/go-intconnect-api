package dashboard_widget

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindBatchById(gormTransaction *gorm.DB, dashboardWidgetIds []uint64) ([]*entity.DashboardWidget, error)
	CreateBatch(gormTransaction *gorm.DB, dashboardWidgets []*entity.DashboardWidget) error
	Update(gormTransaction *gorm.DB, dashboardWidget *entity.DashboardWidget) error
	DeleteBatchById(gormTransaction *gorm.DB, dashboardWidgetIds []uint64) error
	DeleteBatchByCode(gormTransaction *gorm.DB, dashboardWidgetCodes []string) error
}
