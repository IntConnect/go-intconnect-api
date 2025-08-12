package node

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	CreateNode(ginContext *gin.Context)
	DeleteNode(ginContext *gin.Context)
	UpdateNode(ginContext *gin.Context)
}
