package routes

import (
	alarmLog "go-intconnect-api/internal/alarm_log"
	"go-intconnect-api/internal/facility"
	"go-intconnect-api/internal/machine"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	systemSetting "go-intconnect-api/internal/system_setting"

	"github.com/gin-gonic/gin"
)

type PublicRoutes struct {
	systemSettingController systemSetting.Controller
	facilityController      facility.Controller
	machineController       machine.Controller
	alarmLogController      alarmLog.Controller
	mqttBrokerController    mqttBroker.Controller
}

func NewPublicRoutes(routerGroup *gin.RouterGroup, systemSettingController systemSetting.Controller,
	facilityController facility.Controller, machineController machine.Controller, alarmLogController alarmLog.Controller,
	mqttBrokerController mqttBroker.Controller) *PublicRoutes {
	return &PublicRoutes{
		systemSettingController: systemSettingController,
		facilityController:      facilityController,
		machineController:       machineController,
		alarmLogController:      alarmLogController,
		mqttBrokerController:    mqttBrokerController,
	}
}

func (publicRoutes *PublicRoutes) Setup(routerGroup *gin.RouterGroup) {
	publicRouterGroup := routerGroup.Group("/public")

	// Serve static files from "./uploads"
	publicRouterGroup.Static("/uploads", "./uploads")

	publicRouterGroup.GET("/system-settings/:key", publicRoutes.systemSettingController.FindMinimalSystemSettingByKey)
	publicRouterGroup.GET("/facilities", publicRoutes.facilityController.FindMinimalAllFacility)
	publicRouterGroup.GET("/machines/:id", publicRoutes.machineController.FindMinimalMachineById)
	publicRouterGroup.GET("/alarm-logs/machine/:id", publicRoutes.alarmLogController.FindMinimalAlarmLogByMachineId)
	publicRouterGroup.GET("/mqtt-brokers/gateway", publicRoutes.mqttBrokerController.GatewayMqttBroker)

}
