package routes

import (
	"go-intconnect-api/internal/facility"
	systemSetting "go-intconnect-api/internal/system_setting"

	"github.com/gin-gonic/gin"
)

type PublicRoutes struct {
	systemSettingController systemSetting.Controller
	facilityController      facility.Controller
}

func NewPublicRoutes(routerGroup *gin.RouterGroup, systemSettingController systemSetting.Controller,
	facilityController facility.Controller) *PublicRoutes {
	return &PublicRoutes{
		systemSettingController: systemSettingController,
		facilityController:      facilityController,
	}
}

func (publicRoutes *PublicRoutes) Setup(routerGroup *gin.RouterGroup) {
	publicRouterGroup := routerGroup.Group("/public")

	// Serve static files from "./uploads"
	publicRouterGroup.Static("/uploads", "./uploads")

	publicRouterGroup.GET("/system-settings/:key", publicRoutes.systemSettingController.FindMinimalSystemSettingByKey)
	publicRouterGroup.GET("/facilities", publicRoutes.facilityController.FindMinimalAllFacility)
}
