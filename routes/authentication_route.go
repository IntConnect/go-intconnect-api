package routes

import (
	"go-intconnect-api/internal/user"

	"github.com/gin-gonic/gin"
)

type AuthenticationRoutes struct {
	userController user.Controller
}

func NewAuthenticationRoutes(userController user.Controller,
) *AuthenticationRoutes {
	return &AuthenticationRoutes{
		userController: userController,
	}
}

func (authenticationRoutes *AuthenticationRoutes) Setup(routerGroup *gin.RouterGroup) {
	authenticationGroup := routerGroup.Group("authentication")
	authenticationGroup.POST("login", authenticationRoutes.userController.LoginUser)
}
