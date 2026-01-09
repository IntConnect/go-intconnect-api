package telemetry

import "github.com/gin-gonic/gin"

type Controller interface {
	GenerateReport(ginContext *gin.Context)
	IntervalReport(ginContext *gin.Context)
}
