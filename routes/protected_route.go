package routes

import (
	databaseConnection "go-intconnect-api/internal/database_connection"
	"go-intconnect-api/internal/node"
	"go-intconnect-api/internal/pipeline"
	protocolConfiguration "go-intconnect-api/internal/protocol_configuration"
	"go-intconnect-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ProtectedRoutes struct {
	routerGroup                     *gin.RouterGroup
	viperConfig                     *viper.Viper
	userController                  user.Controller
	nodeController                  node.Controller
	pipelineController              pipeline.Controller
	protocolConfigurationController protocolConfiguration.Controller
	databaseConnectionController    databaseConnection.Controller
}

func NewProtectedRoutes(
	routerGroup *gin.RouterGroup,
	viperConfig *viper.Viper,

	userController user.Controller,
	nodeController node.Controller,
	pipelineController pipeline.Controller,
	protocolConfigurationController protocolConfiguration.Controller,
	databaseConnectionController databaseConnection.Controller,
) *ProtectedRoutes {
	//routerGroup.Use(middleware.AuthMiddleware(viperConfig))
	wrapperRouterGroup := routerGroup.Group("/api")
	return &ProtectedRoutes{
		routerGroup: wrapperRouterGroup,
		viperConfig: viperConfig,

		userController:                  userController,
		nodeController:                  nodeController,
		pipelineController:              pipelineController,
		protocolConfigurationController: protocolConfigurationController,
		databaseConnectionController:    databaseConnectionController,
	}
}

func (protectedRoutes *ProtectedRoutes) Setup() {

	userRouterGroup := protectedRoutes.routerGroup.Group("users")
	userRouterGroup.GET("pagination", protectedRoutes.userController.FindAllUserPagination)
	userRouterGroup.GET("", protectedRoutes.userController.FindAllUser)
	userRouterGroup.POST("", protectedRoutes.userController.CreateUser)
	userRouterGroup.PUT("", protectedRoutes.userController.UpdateUser)
	userRouterGroup.DELETE("", protectedRoutes.userController.DeleteUser)

	nodeRouterGroup := protectedRoutes.routerGroup.Group("nodes")
	nodeRouterGroup.GET("pagination", protectedRoutes.nodeController.FindAllPagination)
	nodeRouterGroup.GET("", protectedRoutes.nodeController.FindAll)
	nodeRouterGroup.POST("", protectedRoutes.nodeController.CreateNode)
	nodeRouterGroup.PUT("", protectedRoutes.nodeController.UpdateNode)
	nodeRouterGroup.DELETE("", protectedRoutes.nodeController.DeleteNode)

	pipelineRouterGroup := protectedRoutes.routerGroup.Group("pipelines")
	pipelineRouterGroup.GET("pagination", protectedRoutes.pipelineController.FindAllPagination)
	pipelineRouterGroup.GET("", protectedRoutes.pipelineController.FindAll)
	pipelineRouterGroup.GET("/:id", protectedRoutes.pipelineController.FindById)
	pipelineRouterGroup.POST("", protectedRoutes.pipelineController.CreatePipeline)
	pipelineRouterGroup.GET("/run/:id", protectedRoutes.pipelineController.RunPipeline)
	pipelineRouterGroup.PUT("", protectedRoutes.pipelineController.UpdatePipeline)
	pipelineRouterGroup.DELETE("", protectedRoutes.pipelineController.DeletePipeline)

	protocolConfigurationRouterGroup := protectedRoutes.routerGroup.Group("protocol-configurations")
	protocolConfigurationRouterGroup.GET("pagination", protectedRoutes.protocolConfigurationController.FindAllPagination)
	protocolConfigurationRouterGroup.GET("", protectedRoutes.protocolConfigurationController.FindAll)
	protocolConfigurationRouterGroup.GET("/:id", protectedRoutes.protocolConfigurationController.FindById)
	protocolConfigurationRouterGroup.POST("", protectedRoutes.protocolConfigurationController.CreateProtocolConfiguration)
	protocolConfigurationRouterGroup.PUT("", protectedRoutes.protocolConfigurationController.UpdateProtocolConfiguration)
	protocolConfigurationRouterGroup.DELETE("", protectedRoutes.protocolConfigurationController.DeleteProtocolConfiguration)

	databaseConnectionRouterGroup := protectedRoutes.routerGroup.Group("database-connections")
	databaseConnectionRouterGroup.GET("pagination", protectedRoutes.databaseConnectionController.FindAllPagination)
	databaseConnectionRouterGroup.GET("", protectedRoutes.databaseConnectionController.FindAll)
	databaseConnectionRouterGroup.GET("/:id", protectedRoutes.databaseConnectionController.FindById)
	databaseConnectionRouterGroup.POST("", protectedRoutes.databaseConnectionController.CreateDatabaseConnection)
	databaseConnectionRouterGroup.POST("schema/:id", protectedRoutes.databaseConnectionController.CreateDatabaseSchema)
	databaseConnectionRouterGroup.PUT("", protectedRoutes.databaseConnectionController.UpdateDatabaseConnection)
	databaseConnectionRouterGroup.DELETE("", protectedRoutes.databaseConnectionController.DeleteDatabaseConnection)

}
