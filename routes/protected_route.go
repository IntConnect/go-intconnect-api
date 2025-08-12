package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-intconnect-api/internal/node"
	"go-intconnect-api/internal/user"
	"go-intconnect-api/pkg/middleware"
)

type ProtectedRoutes struct {
	routerGroup    *gin.RouterGroup
	viperConfig    *viper.Viper
	userController user.Controller
	nodeController node.Controller
}

func NewProtectedRoutes(
	routerGroup *gin.RouterGroup,
	viperConfig *viper.Viper,

	userController user.Controller,
	nodeController node.Controller,
) *ProtectedRoutes {
	routerGroup.Use(middleware.AuthMiddleware(viperConfig))

	return &ProtectedRoutes{
		routerGroup: routerGroup,
		viperConfig: viperConfig,

		userController: userController,
		nodeController: nodeController,
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
}
