package permission

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
}
