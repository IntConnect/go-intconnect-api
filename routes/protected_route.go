package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-intconnect-api/internal/user"
	"go-intconnect-api/pkg/middleware"
)

type ProtectedRoutes struct {
	routerGroup    *gin.RouterGroup
	viperConfig    *viper.Viper
	userController user.Controller
}

func NewProtectedRoutes(
	routerGroup *gin.RouterGroup,
	viperConfig *viper.Viper,

	userController user.Controller,
) *ProtectedRoutes {
	routerGroup.Use(middleware.AuthMiddleware(viperConfig))

	return &ProtectedRoutes{
		routerGroup: routerGroup,
		viperConfig: viperConfig,

		userController: userController,
	}
}

func (protectedRoutes *ProtectedRoutes) Setup() {

	userRouterGroup := protectedRoutes.routerGroup.Group("users")
	userRouterGroup.GET("pagination", protectedRoutes.userController.FindAllPagination)
	userRouterGroup.GET("", protectedRoutes.userController.FindAll)
	userRouterGroup.POST("", protectedRoutes.userController.CreateUser)
	userRouterGroup.PUT("", protectedRoutes.userController.UpdateUser)
	userRouterGroup.DELETE("", protectedRoutes.userController.DeleteUser)
}
