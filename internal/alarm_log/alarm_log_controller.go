package alarm_log

import "github.com/gin-gonic/gin"

type Controller interface {
	FindAllAlarmLog(ginContext *gin.Context)
	FindAllAlarmLogPagination(ginContext *gin.Context)
	UpdateAlarmLog(ginContext *gin.Context)
}
