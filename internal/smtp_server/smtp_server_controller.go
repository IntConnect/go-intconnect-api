package smtp_server

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllSmtpServer(ginContext *gin.Context)
	FindAllSmtpServerPagination(ginContext *gin.Context)
	CreateSmtpServer(ginContext *gin.Context)
	DeleteSmtpServer(ginContext *gin.Context)
	UpdateSmtpServer(ginContext *gin.Context)
}
