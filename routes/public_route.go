package routes

import (
	"github.com/gin-gonic/gin"
)

type PublicRoutes struct {
	routerGroup *gin.RouterGroup
}

func NewPublicRoutes(routerGroup *gin.RouterGroup,

) *PublicRoutes {
	return &PublicRoutes{routerGroup: routerGroup.Group("public")}
}

func (publicRoutes *PublicRoutes) Setup() {
	//homeRouterGroup := publicRoutes.routerGroup.Group("/homes/")
}
