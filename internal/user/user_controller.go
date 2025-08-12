package user

import "github.com/gin-gonic/gin"

type Controller interface {
	Login(ginContext *gin.Context)
	FindAll(ginContext *gin.Context)
	FindAllPagination(ginContext *gin.Context)
	CreateUser(ginContext *gin.Context)
	DeleteUser(ginContext *gin.Context)
	UpdateUser(ginContext *gin.Context)
}
