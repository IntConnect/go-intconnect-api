package routes

import (
	"go-intconnect-api/configs"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/breakdown"
	checkSheet "go-intconnect-api/internal/check_sheet"
	checkSheetDocumentTemplate "go-intconnect-api/internal/check_sheet_document_template"
	databaseConnection "go-intconnect-api/internal/database_connection"
	"go-intconnect-api/internal/facility"
	"go-intconnect-api/internal/machine"
	modbusServer "go-intconnect-api/internal/modbus_server"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	mqttTopic "go-intconnect-api/internal/mqtt_topic"
	"go-intconnect-api/internal/node"
	"go-intconnect-api/internal/parameter"
	"go-intconnect-api/internal/permission"
	"go-intconnect-api/internal/pipeline"
	protocolConfiguration "go-intconnect-api/internal/protocol_configuration"
	reportDocumentTemplate "go-intconnect-api/internal/report_document_template"
	"go-intconnect-api/internal/role"
	smtpServer "go-intconnect-api/internal/smtp_server"
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
	nodeController                       node.Controller
	pipelineController                   pipeline.Controller
	protocolConfigurationController      protocolConfiguration.Controller
	databaseConnectionController         databaseConnection.Controller
	facilityController                   facility.Controller
	roleController                       role.Controller
	permissionController                 permission.Controller
	mqttBrokerController                 mqttBroker.Controller
	machineController                    machine.Controller
	parameterController                  parameter.Controller
	mqttTopicController                  mqttTopic.Controller
	auditLogController                   auditLog.Controller
	smtpServerController                 smtpServer.Controller
	modbusServerController               modbusServer.Controller
	reportDocumentTemplateController     reportDocumentTemplate.Controller
	checkSheetDocumentTemplateController checkSheetDocumentTemplate.Controller
	breakdownController                  breakdown.Controller
	systemSettingController              systemSetting.Controller
	telemetryController                  telemetry.Controller
	checkSheetController                 checkSheet.Controller
	roleService                          role.Service
}

func NewProtectedRoutes(
	viperConfig *viper.Viper,
	redisInstance *configs.RedisInstance,

	userController user.Controller,
	nodeController node.Controller,
	pipelineController pipeline.Controller,
	protocolConfigurationController protocolConfiguration.Controller,
	databaseConnectionController databaseConnection.Controller,
	facilityController facility.Controller,
	roleController role.Controller,
	permissionController permission.Controller,
	mqttBrokerController mqttBroker.Controller,
	machineController machine.Controller,
	parameterController parameter.Controller,
	mqttTopicController mqttTopic.Controller, reportDocumentTemplateController reportDocumentTemplate.Controller,
	auditLogController auditLog.Controller,
	smtpServerController smtpServer.Controller,
	modbusServerController modbusServer.Controller,
	checkSheetDocumentTemplateController checkSheetDocumentTemplate.Controller,
	breakdownController breakdown.Controller,
	systemSettingController systemSetting.Controller,
	telemetryController telemetry.Controller,
	checkSheetController checkSheet.Controller,

	roleService role.Service,

) *ProtectedRoutes {
	return &ProtectedRoutes{
		viperConfig: viperConfig,

		userController:                       userController,
		nodeController:                       nodeController,
		pipelineController:                   pipelineController,
		protocolConfigurationController:      protocolConfigurationController,
		databaseConnectionController:         databaseConnectionController,
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
		smtpServerController:                 smtpServerController,
		modbusServerController:               modbusServerController,
		checkSheetDocumentTemplateController: checkSheetDocumentTemplateController,
		breakdownController:                  breakdownController,
		systemSettingController:              systemSettingController,
		telemetryController:                  telemetryController,
		checkSheetController:                 checkSheetController,
	}
}

func (protectedRoutes *ProtectedRoutes) Setup(routerGroup *gin.RouterGroup) {
	routerGroup.Use(middleware.AuthMiddleware(protectedRoutes.viperConfig, protectedRoutes.redisInstance, protectedRoutes.roleService))
	userRouterGroup := routerGroup.Group("users")
	userRouterGroup.GET("pagination", middleware.HasPermission("ROLE_ROLE_USER_VIEW"), protectedRoutes.userController.FindAllUserPagination)
	userRouterGroup.GET("profile", middleware.HasPermission("ROLE_USER_VIEW"), protectedRoutes.userController.FindSelf)
	userRouterGroup.GET("/:id", middleware.HasPermission("ROLE_USER_VIEW"), protectedRoutes.userController.FindById)
	userRouterGroup.POST("", middleware.HasPermission("ROLE_USER_CREATE"), protectedRoutes.userController.CreateUser)
	userRouterGroup.PUT("/profile", protectedRoutes.userController.UpdateProfile)
	userRouterGroup.PUT("/:id", middleware.HasPermission("ROLE_USER_UPDATE"), protectedRoutes.userController.UpdateUser)
	userRouterGroup.DELETE("/:id", middleware.HasPermission("ROLE_USER_DELETE"), protectedRoutes.userController.DeleteUser)
	userRouterGroup.GET("/logout", protectedRoutes.userController.LogoutUser)

	nodeRouterGroup := routerGroup.Group("nodes")
	nodeRouterGroup.GET("pagination", protectedRoutes.nodeController.FindAllPagination)
	nodeRouterGroup.GET("", protectedRoutes.nodeController.FindAll)
	nodeRouterGroup.POST("", protectedRoutes.nodeController.CreateNode)
	nodeRouterGroup.PUT("", protectedRoutes.nodeController.UpdateNode)
	nodeRouterGroup.DELETE("", protectedRoutes.nodeController.DeleteNode)

	pipelineRouterGroup := routerGroup.Group("pipelines")
	pipelineRouterGroup.GET("pagination", protectedRoutes.pipelineController.FindAllPagination)
	pipelineRouterGroup.GET("", protectedRoutes.pipelineController.FindAll, middleware.HasPermission("PIPELINE_VIEW"))
	pipelineRouterGroup.GET("/:id", protectedRoutes.pipelineController.FindById)
	pipelineRouterGroup.POST("", protectedRoutes.pipelineController.CreatePipeline)
	pipelineRouterGroup.GET("/run/:id", protectedRoutes.pipelineController.RunPipeline)
	pipelineRouterGroup.PUT("", protectedRoutes.pipelineController.UpdatePipeline)
	pipelineRouterGroup.DELETE("", protectedRoutes.pipelineController.DeletePipeline)

	protocolConfigurationRouterGroup := routerGroup.Group("protocol-configurations")
	protocolConfigurationRouterGroup.GET("pagination", protectedRoutes.protocolConfigurationController.FindAllPagination)
	protocolConfigurationRouterGroup.GET("", protectedRoutes.protocolConfigurationController.FindAll)
	protocolConfigurationRouterGroup.GET("/:id", protectedRoutes.protocolConfigurationController.FindById)
	protocolConfigurationRouterGroup.POST("", protectedRoutes.protocolConfigurationController.CreateProtocolConfiguration)
	protocolConfigurationRouterGroup.PUT("", protectedRoutes.protocolConfigurationController.UpdateProtocolConfiguration)
	protocolConfigurationRouterGroup.DELETE("", protectedRoutes.protocolConfigurationController.DeleteProtocolConfiguration)

	databaseConnectionRouterGroup := routerGroup.Group("database-connections")
	databaseConnectionRouterGroup.GET("pagination", protectedRoutes.databaseConnectionController.FindAllPagination)
	databaseConnectionRouterGroup.GET("", protectedRoutes.databaseConnectionController.FindAll)
	databaseConnectionRouterGroup.GET("/:id", protectedRoutes.databaseConnectionController.FindById)
	databaseConnectionRouterGroup.POST("", protectedRoutes.databaseConnectionController.CreateDatabaseConnection)
	databaseConnectionRouterGroup.POST("schema/:id", protectedRoutes.databaseConnectionController.CreateDatabaseSchema)
	databaseConnectionRouterGroup.PUT("", protectedRoutes.databaseConnectionController.UpdateDatabaseConnection)
	databaseConnectionRouterGroup.DELETE("", protectedRoutes.databaseConnectionController.DeleteDatabaseConnection)

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

	smtpServerRouterGroup := routerGroup.Group("smtp-servers")
	smtpServerRouterGroup.GET("pagination", protectedRoutes.smtpServerController.FindAllSmtpServerPagination)
	smtpServerRouterGroup.GET("", protectedRoutes.smtpServerController.FindAllSmtpServer)
	smtpServerRouterGroup.POST("", protectedRoutes.smtpServerController.CreateSmtpServer)
	smtpServerRouterGroup.PUT("/:id", protectedRoutes.smtpServerController.UpdateSmtpServer)
	smtpServerRouterGroup.DELETE("/:id", protectedRoutes.smtpServerController.DeleteSmtpServer)

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

	breakdownRouterGroup := routerGroup.Group("breakdowns")
	breakdownRouterGroup.GET("pagination", protectedRoutes.breakdownController.FindAllBreakdownPagination)
	breakdownRouterGroup.GET("", protectedRoutes.breakdownController.FindAllBreakdown)
	breakdownRouterGroup.GET("/:id", protectedRoutes.breakdownController.FindBreakdownById)
	breakdownRouterGroup.POST("", protectedRoutes.breakdownController.CreateBreakdown)
	breakdownRouterGroup.PUT("/:id", protectedRoutes.breakdownController.UpdateBreakdown)
	breakdownRouterGroup.DELETE("/:id", protectedRoutes.breakdownController.DeleteBreakdown)

	systemSettingRouterGroup := routerGroup.Group("system-settings")
	systemSettingRouterGroup.GET("", protectedRoutes.systemSettingController.FindAllSystemSetting)
	systemSettingRouterGroup.GET("/:key", protectedRoutes.systemSettingController.FindSystemSettingByKey)
	systemSettingRouterGroup.POST("", protectedRoutes.systemSettingController.ManageSystemSetting)

	telemetryRouterGroup := routerGroup.Group("telemetries")
	telemetryRouterGroup.POST("/report", protectedRoutes.telemetryController.GenerateReport)

	checkSheetRouterGroup := routerGroup.Group("check-sheets")
	checkSheetRouterGroup.GET("pagination", protectedRoutes.checkSheetController.FindAllCheckSheetPagination)
	checkSheetRouterGroup.GET("", protectedRoutes.checkSheetController.FindAllCheckSheet)
	checkSheetRouterGroup.GET("/:id", protectedRoutes.checkSheetController.FindCheckSheetById)
	checkSheetRouterGroup.POST("", protectedRoutes.checkSheetController.CreateCheckSheet)
	checkSheetRouterGroup.POST("/approval/:id", protectedRoutes.checkSheetController.ApprovalCheckSheet)
	checkSheetRouterGroup.PUT("/:id", protectedRoutes.checkSheetController.UpdateCheckSheet)
	checkSheetRouterGroup.DELETE("/:id", protectedRoutes.checkSheetController.DeleteCheckSheet)

}
