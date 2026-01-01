package dashboard_widget

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (dashboardWidgetRepository *RepositoryImpl) DeleteBatchById(gormTransaction *gorm.DB, dashboardWidgetIds []uint64) error {
	return gormTransaction.Model(&entity.DashboardWidget{}).Where("id IN ?", dashboardWidgetIds).Delete(&entity.DashboardWidget{}).Error
}

func (dashboardWidgetRepository *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, dashboardWidgets []*entity.DashboardWidget) error {
	err := gormTransaction.Create(&dashboardWidgets).Error
	return err
}

func (dashboardWidgetRepository *RepositoryImpl) Update(gormTransaction *gorm.DB, dashboardWidgetEntity *entity.DashboardWidget) error {
	return gormTransaction.Model(dashboardWidgetEntity).Save(dashboardWidgetEntity).Error
}

func (dashboardWidgetRepository *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, dashboardWidgetIds []uint64) ([]*entity.DashboardWidget, error) {
	var dashboardWidgetEntities []*entity.DashboardWidget
	err := gormTransaction.Where("id IN ?", dashboardWidgetIds).Find(&dashboardWidgetEntities).Error
	return dashboardWidgetEntities, err
}
