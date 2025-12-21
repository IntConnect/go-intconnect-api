package injector

import (
	"go-intconnect-api/configs"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/breakdown"
	breakdownResource "go-intconnect-api/internal/breakdown_resource"
	checkSheetDocumentTemplate "go-intconnect-api/internal/check_sheet_document_template"
	databaseConnection "go-intconnect-api/internal/database_connection"
	"go-intconnect-api/internal/facility"
	"go-intconnect-api/internal/machine"
	machineDocument "go-intconnect-api/internal/machine_document"
	modbusServer "go-intconnect-api/internal/modbus_server"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	mqttTopic "go-intconnect-api/internal/mqtt_topic"
	"go-intconnect-api/internal/node"
	"go-intconnect-api/internal/parameter"
	parameterOperation "go-intconnect-api/internal/parameter_operation"
	"go-intconnect-api/internal/permission"
	"go-intconnect-api/internal/pipeline"
	pipelineEdge "go-intconnect-api/internal/pipeline_edge"
	pipelineNode "go-intconnect-api/internal/pipeline_node"
	protocolConfiguration "go-intconnect-api/internal/protocol_configuration"
	reportDocumentTemplate "go-intconnect-api/internal/report_document_template"
	"go-intconnect-api/internal/role"
	smtpServer "go-intconnect-api/internal/smtp_server"
	"go-intconnect-api/internal/storage"
	systemSetting "go-intconnect-api/internal/system_setting"
	"go-intconnect-api/internal/telemetry"
	"go-intconnect-api/internal/user"
	validatorService "go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/middleware"
	"go-intconnect-api/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	universalTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewDatabaseConnection(databaseCredentials *configs.DatabaseCredentials) *gorm.DB {
	databaseConnectionInstance := configs.NewDatabaseConnection(databaseCredentials)
	return databaseConnectionInstance.GetDatabaseConnection()
}

func NewRedisInstance(redisConfig configs.RedisConfig) *configs.RedisInstance {
	redisInstance, err := configs.InitRedisInstance(redisConfig)
	if err != nil {
		panic(err)
	}
	return redisInstance
}

func NewRedisConfig() configs.RedisConfig {
	return configs.RedisConfig{
		IPAddress: "localhost:6379",
		Password:  "",
		Database:  0,
	}
}

func NewValidator(gormDatabase *gorm.DB) (*validator.Validate, universalTranslator.Translator) {
	return configs.InitializeValidator(gormDatabase)
}

// NewViperConfig --- Provider untuk Viper config ---
func NewViperConfig() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.SetConfigFile(".env")
	viperConfig.AddConfigPath(".")
	viperConfig.AutomaticEnv()
	if err := viperConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	return viperConfig
}

// NewDatabaseCredentials --- Provider untuk Database Credentials ---
func NewDatabaseCredentials(viperConfig *viper.Viper) *configs.DatabaseCredentials {
	return &configs.DatabaseCredentials{
		DatabaseHost:     viperConfig.GetString("DATABASE_HOST"),
		DatabasePort:     viperConfig.GetString("DATABASE_PORT"),
		DatabaseName:     viperConfig.GetString("DATABASE_NAME"),
		DatabasePassword: viperConfig.GetString("DATABASE_PASSWORD"),
		DatabaseUsername: viperConfig.GetString("DATABASE_USERNAME"),
	}
}

// NewGinEngine --- Provider untuk Gin Engine ---
func NewGinEngine() (*gin.Engine, *gin.RouterGroup) {
	gin.SetMode(gin.DebugMode)
	ginEngine := gin.Default()
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(exception.Interceptor())
	ginEngineRoot := ginEngine.Group("/")
	ginEngineRoot.Use(middleware.RequestMetaMiddleware())

	return ginEngine, ginEngineRoot
}

func NewStorageManager(viperConfig *viper.Viper) *storage.Manager {
	storageConfig := configs.NewStorageConfig(viperConfig)
	storageManager, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		panic(err)
	}
	return storageManager
}

var CoreModule = fx.Module("coreModule", fx.Provide(
	NewViperConfig,
	NewDatabaseCredentials,
	NewDatabaseConnection,
	NewValidator,
	NewGinEngine,
	NewRedisConfig,
	NewRedisInstance,
	NewStorageManager,
))

var ApplicationRoutesModule = fx.Module("applicationRoutes",
	fx.Provide(
		routes.NewPublicRoutes,
		routes.NewAuthenticationRoutes,
		routes.NewProtectedRoutes,
		func(
			ginEngine *gin.Engine,
			publicRoutes *routes.PublicRoutes,
			authenticationRoutes *routes.AuthenticationRoutes,
			protectedRoutes *routes.ProtectedRoutes,
		) *routes.ApplicationRoutes {
			return routes.NewApplicationRoutes(ginEngine, publicRoutes, authenticationRoutes, protectedRoutes)
		},
	),
	fx.Invoke(func(applicationRoutes *routes.ApplicationRoutes) {
		applicationRoutes.Setup()
	}),
)

var UserModule = fx.Module("userFeature",
	fx.Provide(fx.Annotate(user.NewRepository, fx.As(new(user.Repository)))),
	fx.Provide(fx.Annotate(user.NewService, fx.As(new(user.Service)))),
	fx.Provide(fx.Annotate(user.NewHandler, fx.As(new(user.Controller)))),
)

var NodeModule = fx.Module("nodeFeature",
	fx.Provide(fx.Annotate(node.NewRepository, fx.As(new(node.Repository)))),
	fx.Provide(fx.Annotate(node.NewService, fx.As(new(node.Service)))),
	fx.Provide(fx.Annotate(node.NewHandler, fx.As(new(node.Controller)))),
)

var ValidatorModule = fx.Module("validatorFeature",
	fx.Provide(fx.Annotate(validatorService.NewService, fx.As(new(validatorService.Service)))),
)

var PipelineModule = fx.Module("pipelineFeature",
	fx.Provide(fx.Annotate(pipeline.NewRepository, fx.As(new(pipeline.Repository)))),
	fx.Provide(fx.Annotate(pipeline.NewService, fx.As(new(pipeline.Service)))),
	fx.Provide(fx.Annotate(pipeline.NewHandler, fx.As(new(pipeline.Controller)))),
)

var PipelineConfigurationModule = fx.Module("pipelineConfigurationFeature",
	fx.Provide(fx.Annotate(protocolConfiguration.NewRepository, fx.As(new(protocolConfiguration.Repository)))),
	fx.Provide(fx.Annotate(protocolConfiguration.NewService, fx.As(new(protocolConfiguration.Service)))),
	fx.Provide(fx.Annotate(protocolConfiguration.NewHandler, fx.As(new(protocolConfiguration.Controller)))),
)

var DatabaseConnectionModule = fx.Module("databaseConnectionFeature",
	fx.Provide(fx.Annotate(databaseConnection.NewRepository, fx.As(new(databaseConnection.Repository)))),
	fx.Provide(fx.Annotate(databaseConnection.NewService, fx.As(new(databaseConnection.Service)))),
	fx.Provide(fx.Annotate(databaseConnection.NewHandler, fx.As(new(databaseConnection.Controller)))),
)

var FacilityModule = fx.Module("facilityFeature",
	fx.Provide(fx.Annotate(facility.NewRepository, fx.As(new(facility.Repository)))),
	fx.Provide(fx.Annotate(facility.NewService, fx.As(new(facility.Service)))),
	fx.Provide(fx.Annotate(facility.NewHandler, fx.As(new(facility.Controller)))),
)

var RoleModule = fx.Module("roleFeature",
	fx.Provide(fx.Annotate(role.NewRepository, fx.As(new(role.Repository)))),
	fx.Provide(fx.Annotate(role.NewService, fx.As(new(role.Service)))),
	fx.Provide(fx.Annotate(role.NewHandler, fx.As(new(role.Controller)))),
)

var PermissionModule = fx.Module("permissionFeature",
	fx.Provide(fx.Annotate(permission.NewRepository, fx.As(new(permission.Repository)))),
	fx.Provide(fx.Annotate(permission.NewService, fx.As(new(permission.Service)))),
	fx.Provide(fx.Annotate(permission.NewHandler, fx.As(new(permission.Controller)))),
)

var PipelineNodeModule = fx.Module("pipelineNodeFeature",
	fx.Provide(fx.Annotate(pipelineNode.NewRepository, fx.As(new(pipelineNode.Repository)))),
)

var PipelineEdgeModule = fx.Module("pipelineEdgeFeature",
	fx.Provide(fx.Annotate(pipelineEdge.NewRepository, fx.As(new(pipelineEdge.Repository)))),
)

var MqttBrokerModule = fx.Module("mqttBrokerFeature",
	fx.Provide(fx.Annotate(mqttBroker.NewRepository, fx.As(new(mqttBroker.Repository)))),
	fx.Provide(fx.Annotate(mqttBroker.NewService, fx.As(new(mqttBroker.Service)))),
	fx.Provide(fx.Annotate(mqttBroker.NewHandler, fx.As(new(mqttBroker.Controller)))),
)

var MachineModule = fx.Module("machineFeature",
	fx.Provide(fx.Annotate(machine.NewRepository, fx.As(new(machine.Repository)))),
	fx.Provide(fx.Annotate(machine.NewService, fx.As(new(machine.Service)))),
	fx.Provide(fx.Annotate(machine.NewHandler, fx.As(new(machine.Controller)))),
)

var ParameterModule = fx.Module("parameterFeature",
	fx.Provide(fx.Annotate(parameter.NewRepository, fx.As(new(parameter.Repository)))),
	fx.Provide(fx.Annotate(parameter.NewService, fx.As(new(parameter.Service)))),
	fx.Provide(fx.Annotate(parameter.NewHandler, fx.As(new(parameter.Controller)))),
)

var MqttTopicModule = fx.Module("mqttTopicFeature",
	fx.Provide(fx.Annotate(mqttTopic.NewRepository, fx.As(new(mqttTopic.Repository)))),
	fx.Provide(fx.Annotate(mqttTopic.NewService, fx.As(new(mqttTopic.Service)))),
	fx.Provide(fx.Annotate(mqttTopic.NewHandler, fx.As(new(mqttTopic.Controller)))),
)

var MachineDocumentModule = fx.Module("machineDocumentFeature",
	fx.Provide(fx.Annotate(machineDocument.NewRepository, fx.As(new(machineDocument.Repository)))),
)

var ReportDocumentTemplateModule = fx.Module("reportDocumentTemplateFeature",
	fx.Provide(fx.Annotate(reportDocumentTemplate.NewRepository, fx.As(new(reportDocumentTemplate.Repository)))),
	fx.Provide(fx.Annotate(reportDocumentTemplate.NewService, fx.As(new(reportDocumentTemplate.Service)))),
	fx.Provide(fx.Annotate(reportDocumentTemplate.NewHandler, fx.As(new(reportDocumentTemplate.Controller)))),
)

var AuditLogModule = fx.Module("auditLogFeature",
	fx.Provide(fx.Annotate(auditLog.NewRepository, fx.As(new(auditLog.Repository)))),
	fx.Provide(fx.Annotate(auditLog.NewService, fx.As(new(auditLog.Service)))),
	fx.Provide(fx.Annotate(auditLog.NewHandler, fx.As(new(auditLog.Controller)))),
)

var SmtpServerModule = fx.Module("smtpServerFeature",
	fx.Provide(fx.Annotate(smtpServer.NewRepository, fx.As(new(smtpServer.Repository)))),
	fx.Provide(fx.Annotate(smtpServer.NewService, fx.As(new(smtpServer.Service)))),
	fx.Provide(fx.Annotate(smtpServer.NewHandler, fx.As(new(smtpServer.Controller)))),
)

var ModbusServerModule = fx.Module("modbusServerFeature",
	fx.Provide(fx.Annotate(modbusServer.NewRepository, fx.As(new(modbusServer.Repository)))),
	fx.Provide(fx.Annotate(modbusServer.NewService, fx.As(new(modbusServer.Service)))),
	fx.Provide(fx.Annotate(modbusServer.NewHandler, fx.As(new(modbusServer.Controller)))),
)

var BreakdownModule = fx.Module("breakdownFeature",
	fx.Provide(fx.Annotate(breakdown.NewRepository, fx.As(new(breakdown.Repository)))),
	fx.Provide(fx.Annotate(breakdown.NewService, fx.As(new(breakdown.Service)))),
	fx.Provide(fx.Annotate(breakdown.NewHandler, fx.As(new(breakdown.Controller)))),
)

var CheckSheetDocumentModule = fx.Module("checkSheetDocumentTemplateFeature",
	fx.Provide(fx.Annotate(checkSheetDocumentTemplate.NewRepository, fx.As(new(checkSheetDocumentTemplate.Repository)))),
	fx.Provide(fx.Annotate(checkSheetDocumentTemplate.NewService, fx.As(new(checkSheetDocumentTemplate.Service)))),
	fx.Provide(fx.Annotate(checkSheetDocumentTemplate.NewHandler, fx.As(new(checkSheetDocumentTemplate.Controller)))),
)

var BreakdownResourceModule = fx.Module("breakdownResourceFeature",
	fx.Provide(fx.Annotate(breakdownResource.NewRepository, fx.As(new(breakdownResource.Repository)))),
)

var ParameterOperationModule = fx.Module("parameterOperationFeature",
	fx.Provide(fx.Annotate(parameterOperation.NewRepository, fx.As(new(parameterOperation.Repository)))),
)

var SystemSettingModule = fx.Module("systemSettingFeature",
	fx.Provide(fx.Annotate(systemSetting.NewRepository, fx.As(new(systemSetting.Repository)))),
	fx.Provide(fx.Annotate(systemSetting.NewService, fx.As(new(systemSetting.Service)))),
	fx.Provide(fx.Annotate(systemSetting.NewHandler, fx.As(new(systemSetting.Controller)))),
)

var TelemetryModule = fx.Module("telemetryFeature",
	fx.Provide(fx.Annotate(telemetry.NewRepository, fx.As(new(telemetry.Repository)))),
	fx.Provide(fx.Annotate(telemetry.NewService, fx.As(new(telemetry.Service)))),
	fx.Provide(fx.Annotate(telemetry.NewHandler, fx.As(new(telemetry.Controller)))),
)
