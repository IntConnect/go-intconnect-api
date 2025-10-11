package routes

import (
	"go-intconnect-api/internal/user"

	"github.com/gin-gonic/gin"
)

type AuthenticationRoutes struct {
	routerGroup    *gin.RouterGroup
	userController user.Controller
}

func NewAuthenticationRoutes(routerGroup *gin.RouterGroup, userController user.Controller,
) *AuthenticationRoutes {
	return &AuthenticationRoutes{
		routerGroup:    routerGroup.Group("authentication"),
		userController: userController,
	}
}

func (routerGroup *AuthenticationRoutes) Setup() {
	routerGroup.routerGroup.POST("login", routerGroup.userController.LoginUser)
}
