package routes

import (
	"go-intconnect-api/configs"
	auditLog "go-intconnect-api/internal/audit_log"
	databaseConnection "go-intconnect-api/internal/database_connection"
	"go-intconnect-api/internal/facility"
	"go-intconnect-api/internal/machine"
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
	"go-intconnect-api/internal/user"
	"go-intconnect-api/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ProtectedRoutes struct {
	viperConfig                      *viper.Viper
	userController                   user.Controller
	nodeController                   node.Controller
	pipelineController               pipeline.Controller
	protocolConfigurationController  protocolConfiguration.Controller
	databaseConnectionController     databaseConnection.Controller
	facilityController               facility.Controller
	roleController                   role.Controller
	roleService                      role.Service
	permissionController             permission.Controller
	mqttBrokerController             mqttBroker.Controller
	machineController                machine.Controller
	parameterController              parameter.Controller
	mqttTopicController              mqttTopic.Controller
	reportDocumentTemplateController reportDocumentTemplate.Controller
	auditLogController               auditLog.Controller
	redisInstance                    *configs.RedisInstance
	smtpServerController             smtpServer.Controller
}

func NewProtectedRoutes(
	viperConfig *viper.Viper,

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
	redisInstance *configs.RedisInstance,
	roleService role.Service,
	smtpServerController smtpServer.Controller,
) *ProtectedRoutes {
	return &ProtectedRoutes{
		viperConfig: viperConfig,

		userController:                   userController,
		nodeController:                   nodeController,
		pipelineController:               pipelineController,
		protocolConfigurationController:  protocolConfigurationController,
		databaseConnectionController:     databaseConnectionController,
		facilityController:               facilityController,
		roleController:                   roleController,
		permissionController:             permissionController,
		mqttBrokerController:             mqttBrokerController,
		machineController:                machineController,
		parameterController:              parameterController,
		mqttTopicController:              mqttTopicController,
		reportDocumentTemplateController: reportDocumentTemplateController,
		auditLogController:               auditLogController,
		redisInstance:                    redisInstance,
		roleService:                      roleService,
		smtpServerController:             smtpServerController,
	}
}

func (protectedRoutes *ProtectedRoutes) Setup(routerGroup *gin.RouterGroup) {
	routerGroup.Use(middleware.AuthMiddleware(protectedRoutes.viperConfig, protectedRoutes.redisInstance, protectedRoutes.roleService))
	userRouterGroup := routerGroup.Group("users")
	userRouterGroup.GET("pagination", protectedRoutes.userController.FindAllUserPagination, middleware.HasPermission("USER_VIEW"))
	userRouterGroup.GET("", protectedRoutes.userController.FindAllUser, middleware.HasPermission("USER_VIEW"))
	userRouterGroup.GET("/:id", protectedRoutes.userController.FindById, middleware.HasPermission("USER_VIEW"))
	userRouterGroup.POST("", protectedRoutes.userController.CreateUser, middleware.HasPermission("USER_CREATE"))
	userRouterGroup.PUT("/:id", protectedRoutes.userController.UpdateUser, middleware.HasPermission("USER_UPDATE"))
	userRouterGroup.DELETE("/:id", protectedRoutes.userController.DeleteUser, middleware.HasPermission("USER_DELETE"))

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
	facilityRouterGroup.GET("pagination", protectedRoutes.facilityController.FindAllPagination)
	facilityRouterGroup.GET("", protectedRoutes.facilityController.FindAll)
	facilityRouterGroup.POST("", protectedRoutes.facilityController.CreateFacility)
	facilityRouterGroup.PUT("/:id", protectedRoutes.facilityController.UpdateFacility)
	facilityRouterGroup.DELETE("/:id", protectedRoutes.facilityController.DeleteFacility)

	roleRouterGroup := routerGroup.Group("roles")
	roleRouterGroup.GET("", protectedRoutes.roleController.FindAllRole)
	roleRouterGroup.POST("", protectedRoutes.roleController.CreateRole)
	roleRouterGroup.PUT("", protectedRoutes.roleController.UpdateRole)
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
	machineRouterGroup.POST("", protectedRoutes.machineController.CreateMachine)
	machineRouterGroup.PUT("", protectedRoutes.machineController.UpdateMachine)
	machineRouterGroup.DELETE("", protectedRoutes.machineController.DeleteMachine)

	parameterRouterGroup := routerGroup.Group("parameters")
	parameterRouterGroup.GET("pagination", protectedRoutes.parameterController.FindAllParameterPagination)
	parameterRouterGroup.GET("", protectedRoutes.parameterController.FindAllParameter)
	parameterRouterGroup.GET("/create", protectedRoutes.parameterController.FindDependencyParameter)
	parameterRouterGroup.POST("", protectedRoutes.parameterController.CreateParameter)
	parameterRouterGroup.PUT("", protectedRoutes.parameterController.UpdateParameter)
	parameterRouterGroup.DELETE("", protectedRoutes.parameterController.DeleteParameter)

	mqttTopicRouterGroup := routerGroup.Group("mqtt-topics")
	mqttTopicRouterGroup.GET("pagination", protectedRoutes.mqttTopicController.FindAllMqttTopicPagination)
	mqttTopicRouterGroup.GET("", protectedRoutes.mqttTopicController.FindAllMqttTopic)
	mqttTopicRouterGroup.GET("create", protectedRoutes.mqttTopicController.FindDependencyMqttTopic)
	mqttTopicRouterGroup.POST("", protectedRoutes.mqttTopicController.CreateMqttTopic)
	mqttTopicRouterGroup.PUT("", protectedRoutes.mqttTopicController.UpdateMqttTopic)
	mqttTopicRouterGroup.DELETE("", protectedRoutes.mqttTopicController.DeleteMqttTopic)

	reportDocumentTemplateRouterGroup := routerGroup.Group("report-document-templates")
	reportDocumentTemplateRouterGroup.GET("pagination", protectedRoutes.reportDocumentTemplateController.FindAllReportDocumentTemplatePagination)
	reportDocumentTemplateRouterGroup.GET("", protectedRoutes.reportDocumentTemplateController.FindAllReportDocumentTemplate)
	reportDocumentTemplateRouterGroup.POST("", protectedRoutes.reportDocumentTemplateController.CreateReportDocumentTemplate)
	reportDocumentTemplateRouterGroup.PUT("", protectedRoutes.reportDocumentTemplateController.UpdateReportDocumentTemplate)
	reportDocumentTemplateRouterGroup.DELETE("", protectedRoutes.reportDocumentTemplateController.DeleteReportDocumentTemplate)

	auditLogRouterGroup := routerGroup.Group("audit-logs")
	auditLogRouterGroup.GET("pagination", protectedRoutes.auditLogController.FindAllAuditLogPagination)
	auditLogRouterGroup.GET("", protectedRoutes.auditLogController.FindAllAuditLog)

	smtpServerRouterGroup := routerGroup.Group("smtp-servers")
	smtpServerRouterGroup.GET("pagination", protectedRoutes.smtpServerController.FindAllSmtpServerPagination)
	smtpServerRouterGroup.GET("", protectedRoutes.smtpServerController.FindAllSmtpServer)
	smtpServerRouterGroup.POST("", protectedRoutes.smtpServerController.CreateSmtpServer)
	smtpServerRouterGroup.PUT("/:id", protectedRoutes.smtpServerController.UpdateSmtpServer)
	smtpServerRouterGroup.DELETE("", protectedRoutes.smtpServerController.DeleteSmtpServer)

}
