package database_connection

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	FindById(ginContext *gin.Context)
	CreateDatabaseConnection(ginContext *gin.Context)
	DeleteDatabaseConnection(ginContext *gin.Context)
	UpdateDatabaseConnection(ginContext *gin.Context)
	CreateDatabaseSchema(context *gin.Context)
}
