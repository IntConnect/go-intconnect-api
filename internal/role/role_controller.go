package role

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllRole(ginContext *gin.Context)
	CreateRole(ginContext *gin.Context)
	DeleteRole(ginContext *gin.Context)
	UpdateRole(ginContext *gin.Context)
}
