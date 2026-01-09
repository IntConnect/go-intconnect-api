package injector

import (
	"go-intconnect-api/configs"
	alarmLog "go-intconnect-api/internal/alarm_log"
	auditLog "go-intconnect-api/internal/audit_log"
	checkSheet "go-intconnect-api/internal/check_sheet"
	checkSheetDocumentTemplate "go-intconnect-api/internal/check_sheet_document_template"
	checkSheetValue "go-intconnect-api/internal/check_sheet_value"
	dashboardWidget "go-intconnect-api/internal/dashboard_widget"
	"go-intconnect-api/internal/facility"
	"go-intconnect-api/internal/machine"
	machineDocument "go-intconnect-api/internal/machine_document"
	modbusServer "go-intconnect-api/internal/modbus_server"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	mqttTopic "go-intconnect-api/internal/mqtt_topic"
	"go-intconnect-api/internal/parameter"
	parameterOperation "go-intconnect-api/internal/parameter_operation"
	"go-intconnect-api/internal/permission"
	processedParameterSequence "go-intconnect-api/internal/processed_parameter_sequence"
	"go-intconnect-api/internal/register"
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

var ValidatorModule = fx.Module("validatorFeature",
	fx.Provide(fx.Annotate(validatorService.NewService, fx.As(new(validatorService.Service)))),
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

var CheckSheetDocumentModule = fx.Module("checkSheetDocumentTemplateFeature",
	fx.Provide(fx.Annotate(checkSheetDocumentTemplate.NewRepository, fx.As(new(checkSheetDocumentTemplate.Repository)))),
	fx.Provide(fx.Annotate(checkSheetDocumentTemplate.NewService, fx.As(new(checkSheetDocumentTemplate.Service)))),
	fx.Provide(fx.Annotate(checkSheetDocumentTemplate.NewHandler, fx.As(new(checkSheetDocumentTemplate.Controller)))),
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

var CheckSheetModule = fx.Module("checkSheetFeature",
	fx.Provide(fx.Annotate(checkSheet.NewRepository, fx.As(new(checkSheet.Repository)))),
	fx.Provide(fx.Annotate(checkSheet.NewService, fx.As(new(checkSheet.Service)))),
	fx.Provide(fx.Annotate(checkSheet.NewHandler, fx.As(new(checkSheet.Controller)))),
)

var CheckSheetValueModule = fx.Module("checkSheetValueFeature",
	fx.Provide(fx.Annotate(checkSheetValue.NewRepository, fx.As(new(checkSheetValue.Repository)))),
)

var DashboardWidgetModule = fx.Module("dashboardWidgetFeature",
	fx.Provide(fx.Annotate(dashboardWidget.NewRepository, fx.As(new(dashboardWidget.Repository)))),
)

var RegisterModule = fx.Module("registerFeature",
	fx.Provide(fx.Annotate(register.NewRepository, fx.As(new(register.Repository)))),
	fx.Provide(fx.Annotate(register.NewService, fx.As(new(register.Service)))),
	fx.Provide(fx.Annotate(register.NewHandler, fx.As(new(register.Controller)))),
)

var AlarmLogModule = fx.Module("alarmLogFeature",
	fx.Provide(fx.Annotate(alarmLog.NewRepository, fx.As(new(alarmLog.Repository)))),
	fx.Provide(fx.Annotate(alarmLog.NewService, fx.As(new(alarmLog.Service)))),
	fx.Provide(fx.Annotate(alarmLog.NewHandler, fx.As(new(alarmLog.Controller)))),
)

var ProcessedParameterSequenceModule = fx.Module("processedParameterSequenceFeature",
	fx.Provide(fx.Annotate(processedParameterSequence.NewRepository, fx.As(new(processedParameterSequence.Repository)))),
)
