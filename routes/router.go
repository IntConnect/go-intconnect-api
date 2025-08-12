package routes

import (
	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	RouterGroup *gin.RouterGroup
}

// ApplicationRoutes holds all route groups
type ApplicationRoutes struct {
	PublicRoutes         *PublicRoutes
	AuthenticationRoutes *AuthenticationRoutes
	ProtectedRoutes      *ProtectedRoutes
	// Add other route groups here
}
