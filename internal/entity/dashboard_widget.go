package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type DashboardWidget struct {
	Id        uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId uint64                 `gorm:"column:machine_id;"`
	Code      string                 `gorm:"column:code;"`
	Layout    map[string]interface{} `gorm:"-:all"`
	Config    map[string]interface{} `gorm:"-:all"`
	LayoutRaw []byte                 `gorm:"column:layout;type:jsonb"`
	ConfigRaw []byte                 `gorm:"column:config;type:jsonb"`
	Machine   Machine                `gorm:"foreignKey:MachineId;references:Id"`
}

func (dashboardWidgetEntity *DashboardWidget) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(dashboardWidgetEntity.LayoutRaw) > 0 {
		err = json.Unmarshal(dashboardWidgetEntity.LayoutRaw, &dashboardWidgetEntity.Layout)
	}

	if len(dashboardWidgetEntity.ConfigRaw) > 0 {
		err = json.Unmarshal(dashboardWidgetEntity.ConfigRaw, &dashboardWidgetEntity.Config)
	}
	return
}

func (dashboardWidgetEntity *DashboardWidget) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if dashboardWidgetEntity.Layout != nil {
		dashboardWidgetEntity.LayoutRaw, err = json.Marshal(dashboardWidgetEntity.Layout)
	}
	if dashboardWidgetEntity.Config != nil {
		dashboardWidgetEntity.ConfigRaw, err = json.Marshal(dashboardWidgetEntity.Config)
	}
	return
}
