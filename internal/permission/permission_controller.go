package permission

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllPermission(ginContext *gin.Context)
	FindAllPermissionPagination(ginContext *gin.Context)
}
