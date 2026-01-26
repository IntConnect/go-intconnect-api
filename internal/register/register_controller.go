package register

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllRegister(ginContext *gin.Context)
	FindAllRegisterPagination(ginContext *gin.Context)
	FindRegisterById(ginContext *gin.Context)
	FindRegisterDependency(ginContext *gin.Context)
	CreateRegister(ginContext *gin.Context)
	DeleteRegister(ginContext *gin.Context)
	UpdateRegister(ginContext *gin.Context)
	UpdateRegisterValue(ginContext *gin.Context)
}
