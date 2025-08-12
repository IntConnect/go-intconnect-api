package routes

import (
	"github.com/gin-gonic/gin"
	"go-intconnect-api/internal/home"
)

type PublicRoutes struct {
	routerGroup    *gin.RouterGroup
	homeController home.Controller
}

func NewPublicRoutes(routerGroup *gin.RouterGroup,
	homeController home.Controller,

) *PublicRoutes {
	return &PublicRoutes{routerGroup: routerGroup.Group("public"),
		homeController: homeController,
	}
}

func (publicRoutes *PublicRoutes) Setup() {
	homeRouterGroup := publicRoutes.routerGroup.Group("/homes/")
	homeRouterGroup.GET("", publicRoutes.homeController.Ping)
	homeRouterGroup.GET("dashboard", publicRoutes.homeController.Dashboard)
}
