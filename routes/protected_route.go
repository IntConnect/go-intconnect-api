package routes

import (
	"go-intconnect-api/configs"
	alarmLog "go-intconnect-api/internal/alarm_log"
	auditLog "go-intconnect-api/internal/audit_log"
	checkSheet "go-intconnect-api/internal/check_sheet"
	checkSheetDocumentTemplate "go-intconnect-api/internal/check_sheet_document_template"
	"go-intconnect-api/internal/facility"
	"go-intconnect-api/internal/machine"
	modbusServer "go-intconnect-api/internal/modbus_server"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	mqttTopic "go-intconnect-api/internal/mqtt_topic"
	"go-intconnect-api/internal/parameter"
	"go-intconnect-api/internal/permission"
	"go-intconnect-api/internal/register"
	reportDocumentTemplate "go-intconnect-api/internal/report_document_template"
	"go-intconnect-api/internal/role"
	systemSetting "go-intconnect-api/internal/system_setting"
	"go-intconnect-api/internal/telemetry"
	"go-intconnect-api/internal/user"
	"go-intconnect-api/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ProtectedRoutes struct {
	viperConfig                          *viper.Viper
	redisInstance                        *configs.RedisInstance
	userController                       user.Controller
	facilityController                   facility.Controller
	roleController                       role.Controller
	permissionController                 permission.Controller
	mqttBrokerController                 mqttBroker.Controller
	machineController                    machine.Controller
	parameterController                  parameter.Controller
	mqttTopicController                  mqttTopic.Controller
	auditLogController                   auditLog.Controller
	modbusServerController               modbusServer.Controller
	reportDocumentTemplateController     reportDocumentTemplate.Controller
	checkSheetDocumentTemplateController checkSheetDocumentTemplate.Controller
	systemSettingController              systemSetting.Controller
	telemetryController                  telemetry.Controller
	checkSheetController                 checkSheet.Controller
	registerController                   register.Controller
	alarmLogController                   alarmLog.Controller

	roleService role.Service
}

func NewProtectedRoutes(
	viperConfig *viper.Viper,
	redisInstance *configs.RedisInstance,

	userController user.Controller,
	facilityController facility.Controller,
	roleController role.Controller,
	permissionController permission.Controller,
	mqttBrokerController mqttBroker.Controller,
	machineController machine.Controller,
	parameterController parameter.Controller,
	mqttTopicController mqttTopic.Controller, reportDocumentTemplateController reportDocumentTemplate.Controller,
	auditLogController auditLog.Controller,
	modbusServerController modbusServer.Controller,
	checkSheetDocumentTemplateController checkSheetDocumentTemplate.Controller,
	systemSettingController systemSetting.Controller,
	telemetryController telemetry.Controller,
	checkSheetController checkSheet.Controller,
	registerController register.Controller,
	alarmLogController alarmLog.Controller,
	roleService role.Service,

) *ProtectedRoutes {
	return &ProtectedRoutes{
		viperConfig: viperConfig,

		userController:                       userController,
		facilityController:                   facilityController,
		roleController:                       roleController,
		permissionController:                 permissionController,
		mqttBrokerController:                 mqttBrokerController,
		machineController:                    machineController,
		parameterController:                  parameterController,
		mqttTopicController:                  mqttTopicController,
		reportDocumentTemplateController:     reportDocumentTemplateController,
		auditLogController:                   auditLogController,
		redisInstance:                        redisInstance,
		roleService:                          roleService,
		modbusServerController:               modbusServerController,
		checkSheetDocumentTemplateController: checkSheetDocumentTemplateController,
		systemSettingController:              systemSettingController,
		telemetryController:                  telemetryController,
		checkSheetController:                 checkSheetController,
		registerController:                   registerController,
		alarmLogController:                   alarmLogController,
	}
}

func (protectedRoutes *ProtectedRoutes) Setup(routerGroup *gin.RouterGroup) {
	routerGroup.Use(middleware.AuthMiddleware(protectedRoutes.viperConfig, protectedRoutes.redisInstance, protectedRoutes.roleService))
	userRouterGroup := routerGroup.Group("users")
	userRouterGroup.GET("pagination", middleware.HasPermission("ROLE_USER_VIEW"), protectedRoutes.userController.FindAllUserPagination)
	userRouterGroup.GET("profile", protectedRoutes.userController.FindSelf)
	userRouterGroup.GET("/:id", middleware.HasPermission("ROLE_USER_VIEW"), protectedRoutes.userController.FindById)
	userRouterGroup.POST("", middleware.HasPermission("ROLE_USER_CREATE"), protectedRoutes.userController.CreateUser)
	userRouterGroup.PUT("/profile", protectedRoutes.userController.UpdateProfile)
	userRouterGroup.PUT("/:id", middleware.HasPermission("ROLE_USER_EDIT"), protectedRoutes.userController.UpdateUser)
	userRouterGroup.DELETE("/:id", middleware.HasPermission("ROLE_USER_DELETE"), protectedRoutes.userController.DeleteUser)
	userRouterGroup.GET("/logout", protectedRoutes.userController.LogoutUser)

	facilityRouterGroup := routerGroup.Group("facilities")
	facilityRouterGroup.GET("pagination", protectedRoutes.facilityController.FindAllFacilityPagination)
	facilityRouterGroup.GET("", protectedRoutes.facilityController.FindAllFacility)
	facilityRouterGroup.GET("/:id", protectedRoutes.facilityController.FindFacilityById)
	facilityRouterGroup.POST("", protectedRoutes.facilityController.CreateFacility)
	facilityRouterGroup.PUT("/:id", protectedRoutes.facilityController.UpdateFacility)
	facilityRouterGroup.DELETE("/:id", protectedRoutes.facilityController.DeleteFacility)

	roleRouterGroup := routerGroup.Group("roles")
	roleRouterGroup.GET("", protectedRoutes.roleController.FindAllRole)
	roleRouterGroup.POST("", protectedRoutes.roleController.CreateRole)
	roleRouterGroup.PUT("/:id", protectedRoutes.roleController.UpdateRole)
	roleRouterGroup.DELETE("/:id", protectedRoutes.roleController.DeleteRole)

	permissionRouterGroup := routerGroup.Group("permissions")
	permissionRouterGroup.GET("pagination", protectedRoutes.permissionController.FindAllPermissionPagination)
	permissionRouterGroup.GET("", protectedRoutes.permissionController.FindAllPermission)

	mqttBrokerRouterGroup := routerGroup.Group("mqtt-brokers")
	mqttBrokerRouterGroup.GET("pagination", protectedRoutes.mqttBrokerController.FindAllMqttBrokerPagination)
	mqttBrokerRouterGroup.GET("", protectedRoutes.mqttBrokerController.FindAllMqttBroker)
	mqttBrokerRouterGroup.POST("", protectedRoutes.mqttBrokerController.CreateMqttBroker)
	mqttBrokerRouterGroup.PUT("/:id", protectedRoutes.mqttBrokerController.UpdateMqttBroker)
	mqttBrokerRouterGroup.DELETE("/:id", protectedRoutes.mqttBrokerController.DeleteMqttBroker)

	machineRouterGroup := routerGroup.Group("machines")
	machineRouterGroup.GET("pagination", protectedRoutes.machineController.FindAllMachinePagination)
	machineRouterGroup.GET("", protectedRoutes.machineController.FindAllMachine)
	machineRouterGroup.GET("/:id", protectedRoutes.machineController.FindMachineById)
	machineRouterGroup.GET("/facilities/:id", protectedRoutes.machineController.FindMachineByFacilityId)
	machineRouterGroup.POST("", protectedRoutes.machineController.CreateMachine)
	machineRouterGroup.POST("/dashboard/:id", protectedRoutes.machineController.ManageDashboard)
	machineRouterGroup.PUT("/:id", protectedRoutes.machineController.UpdateMachine)
	machineRouterGroup.DELETE("/:id", protectedRoutes.machineController.DeleteMachine)

	parameterRouterGroup := routerGroup.Group("parameters")
	parameterRouterGroup.GET("pagination", protectedRoutes.parameterController.FindAllParameterPagination)
	parameterRouterGroup.GET("", protectedRoutes.parameterController.FindAllParameter)
	parameterRouterGroup.GET("/:id", protectedRoutes.parameterController.FindByIdParameter)
	parameterRouterGroup.GET("/create", protectedRoutes.parameterController.FindDependencyParameter)
	parameterRouterGroup.POST("", protectedRoutes.parameterController.CreateParameter)
	parameterRouterGroup.PUT("/:id", protectedRoutes.parameterController.UpdateParameter)
	parameterRouterGroup.PUT("operation/:id", protectedRoutes.parameterController.UpdateParameterOperation)
	parameterRouterGroup.DELETE("", protectedRoutes.parameterController.DeleteParameter)

	mqttTopicRouterGroup := routerGroup.Group("mqtt-topics")
	mqttTopicRouterGroup.GET("pagination", protectedRoutes.mqttTopicController.FindAllMqttTopicPagination)
	mqttTopicRouterGroup.GET("", protectedRoutes.mqttTopicController.FindAllMqttTopic)
	mqttTopicRouterGroup.GET("create", protectedRoutes.mqttTopicController.FindDependencyMqttTopic)
	mqttTopicRouterGroup.POST("", protectedRoutes.mqttTopicController.CreateMqttTopic)
	mqttTopicRouterGroup.PUT("/:id", protectedRoutes.mqttTopicController.UpdateMqttTopic)
	mqttTopicRouterGroup.DELETE("", protectedRoutes.mqttTopicController.DeleteMqttTopic)

	reportDocumentTemplateRouterGroup := routerGroup.Group("report-document-templates")
	reportDocumentTemplateRouterGroup.GET("pagination", protectedRoutes.reportDocumentTemplateController.FindAllReportDocumentTemplatePagination)
	reportDocumentTemplateRouterGroup.GET("", protectedRoutes.reportDocumentTemplateController.FindAllReportDocumentTemplate)
	reportDocumentTemplateRouterGroup.POST("", protectedRoutes.reportDocumentTemplateController.CreateReportDocumentTemplate)
	reportDocumentTemplateRouterGroup.PUT("/:id", protectedRoutes.reportDocumentTemplateController.UpdateReportDocumentTemplate)
	reportDocumentTemplateRouterGroup.DELETE("/:id", protectedRoutes.reportDocumentTemplateController.DeleteReportDocumentTemplate)

	auditLogRouterGroup := routerGroup.Group("audit-logs")
	auditLogRouterGroup.GET("pagination", protectedRoutes.auditLogController.FindAllAuditLogPagination)
	auditLogRouterGroup.GET("", protectedRoutes.auditLogController.FindAllAuditLog)

	modbusServerRouterGroup := routerGroup.Group("modbus-servers")
	modbusServerRouterGroup.GET("pagination", protectedRoutes.modbusServerController.FindAllModbusServerPagination)
	modbusServerRouterGroup.GET("", protectedRoutes.modbusServerController.FindAllModbusServer)
	modbusServerRouterGroup.POST("", protectedRoutes.modbusServerController.CreateModbusServer)
	modbusServerRouterGroup.PUT("/:id", protectedRoutes.modbusServerController.UpdateModbusServer)
	modbusServerRouterGroup.DELETE("/:id", protectedRoutes.modbusServerController.DeleteModbusServer)

	checkSheetDocumentTemplateRouterGroup := routerGroup.Group("check-sheet-document-templates")
	checkSheetDocumentTemplateRouterGroup.GET("pagination", protectedRoutes.checkSheetDocumentTemplateController.FindAllCheckSheetDocumentTemplatePagination)
	checkSheetDocumentTemplateRouterGroup.GET("", protectedRoutes.checkSheetDocumentTemplateController.FindAllCheckSheetDocumentTemplate)
	checkSheetDocumentTemplateRouterGroup.POST("", protectedRoutes.checkSheetDocumentTemplateController.CreateCheckSheetDocumentTemplate)
	checkSheetDocumentTemplateRouterGroup.PUT("/:id", protectedRoutes.checkSheetDocumentTemplateController.UpdateCheckSheetDocumentTemplate)
	checkSheetDocumentTemplateRouterGroup.DELETE("/:id", protectedRoutes.checkSheetDocumentTemplateController.DeleteCheckSheetDocumentTemplate)

	systemSettingRouterGroup := routerGroup.Group("system-settings")
	systemSettingRouterGroup.GET("", protectedRoutes.systemSettingController.FindAllSystemSetting)
	systemSettingRouterGroup.GET("/:key", protectedRoutes.systemSettingController.FindSystemSettingByKey)
	systemSettingRouterGroup.POST("", protectedRoutes.systemSettingController.ManageSystemSetting)

	telemetryRouterGroup := routerGroup.Group("telemetries")
	telemetryRouterGroup.POST("/report", protectedRoutes.telemetryController.GenerateReport)
	telemetryRouterGroup.POST("/interval", protectedRoutes.telemetryController.IntervalReport)

	checkSheetRouterGroup := routerGroup.Group("check-sheets")
	checkSheetRouterGroup.GET("pagination", protectedRoutes.checkSheetController.FindAllCheckSheetPagination)
	checkSheetRouterGroup.GET("", protectedRoutes.checkSheetController.FindAllCheckSheet)
	checkSheetRouterGroup.GET("/:id", protectedRoutes.checkSheetController.FindCheckSheetById)
	checkSheetRouterGroup.POST("", protectedRoutes.checkSheetController.CreateCheckSheet)
	checkSheetRouterGroup.POST("/approval/:id", protectedRoutes.checkSheetController.ApprovalCheckSheet)
	checkSheetRouterGroup.PUT("/:id", protectedRoutes.checkSheetController.UpdateCheckSheet)
	checkSheetRouterGroup.DELETE("/:id", protectedRoutes.checkSheetController.DeleteCheckSheet)

	registerRouterGroup := routerGroup.Group("registers")
	registerRouterGroup.GET("", protectedRoutes.registerController.FindAllRegister)
	registerRouterGroup.GET("/:id", protectedRoutes.registerController.FindRegisterById)
	registerRouterGroup.GET("/pagination", protectedRoutes.registerController.FindAllRegisterPagination)
	registerRouterGroup.GET("/dependency", protectedRoutes.registerController.FindRegisterDependency)
	registerRouterGroup.POST("", protectedRoutes.registerController.CreateRegister)
	registerRouterGroup.PUT("/:id", protectedRoutes.registerController.UpdateRegister)
	registerRouterGroup.PUT("/value/:id", protectedRoutes.registerController.UpdateRegisterValue)
	registerRouterGroup.DELETE("/:id", protectedRoutes.registerController.DeleteRegister)

	alarmLogRouterGroup := routerGroup.Group("alarm-logs")
	alarmLogRouterGroup.GET("pagination", protectedRoutes.alarmLogController.FindAllAlarmLogPagination)
	alarmLogRouterGroup.GET("", protectedRoutes.alarmLogController.FindAllAlarmLog)
	alarmLogRouterGroup.PUT("/:id", protectedRoutes.alarmLogController.UpdateAlarmLog)
	alarmLogRouterGroup.GET("/machine/:id", protectedRoutes.alarmLogController.FindAlarmLogByMachineId)

}
