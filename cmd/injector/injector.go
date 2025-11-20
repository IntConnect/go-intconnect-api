package injector

import (
	"go-intconnect-api/configs"
	databaseConnection "go-intconnect-api/internal/database_connection"
	"go-intconnect-api/internal/facility"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
	"go-intconnect-api/internal/node"
	"go-intconnect-api/internal/permission"
	"go-intconnect-api/internal/pipeline"
	pipelineEdge "go-intconnect-api/internal/pipeline_edge"
	pipelineNode "go-intconnect-api/internal/pipeline_node"
	protocolConfiguration "go-intconnect-api/internal/protocol_configuration"
	"go-intconnect-api/internal/role"
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

// --- Provider untuk Viper config ---
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

// --- Provider untuk Database Credentials ---
func NewDatabaseCredentials(viperConfig *viper.Viper) *configs.DatabaseCredentials {
	return &configs.DatabaseCredentials{
		DatabaseHost:     viperConfig.GetString("DATABASE_HOST"),
		DatabasePort:     viperConfig.GetString("DATABASE_PORT"),
		DatabaseName:     viperConfig.GetString("DATABASE_NAME"),
		DatabasePassword: viperConfig.GetString("DATABASE_PASSWORD"),
		DatabaseUsername: viperConfig.GetString("DATABASE_USERNAME"),
	}
}

// --- Provider untuk Gin Engine ---
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

var CoreModule = fx.Module("coreModule", fx.Provide(
	NewViperConfig,
	NewDatabaseCredentials,
	NewDatabaseConnection,
	NewValidator,
	NewGinEngine,
	NewRedisConfig,
	NewRedisInstance,
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

var PipelineEdgeModule = fx.Module("pipelineFeature",
	fx.Provide(fx.Annotate(pipelineEdge.NewRepository, fx.As(new(pipelineEdge.Repository)))),
)

var LoggerModule = fx.Module("loggerFeature", fx.Provide(configs.GetLogger))

var MqttBrokerModule = fx.Module("mqttBrokerFeature",
	fx.Provide(fx.Annotate(mqttBroker.NewRepository, fx.As(new(mqttBroker.Repository)))),
	fx.Provide(fx.Annotate(mqttBroker.NewService, fx.As(new(mqttBroker.Service)))),
	fx.Provide(fx.Annotate(mqttBroker.NewHandler, fx.As(new(mqttBroker.Controller)))),
)
