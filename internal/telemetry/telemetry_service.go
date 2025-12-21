package telemetry

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GenerateReport(ginContext *gin.Context, telemetryReportFilterRequest *model.TelemetryReportFilterRequest) []*model.TelemetryGrouped
}
