package telemetry

import (
	"fmt"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	telemetryService Service
	viperConfig      *viper.Viper
}

func NewHandler(telemetryService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		telemetryService: telemetryService,
		viperConfig:      viperConfig,
	}
}

func (telemetryHandler *Handler) GenerateReport(ginContext *gin.Context) {
	var telemetryReportFilterRequest model.TelemetryReportFilterRequest
	err := ginContext.ShouldBindBodyWithJSON(&telemetryReportFilterRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	telemetryResponses := telemetryHandler.telemetryService.GenerateReport(ginContext, &telemetryReportFilterRequest)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Parameters has been fetched", telemetryResponses))
}

func (telemetryHandler *Handler) GenerateXLSX(ginContext *gin.Context) {
	var telemetryReportFilterRequest model.TelemetryReportFilterRequest
	err := ginContext.ShouldBindBodyWithJSON(&telemetryReportFilterRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	excelFile, err := telemetryHandler.telemetryService.GenerateXLSX(ginContext, &telemetryReportFilterRequest)
	filename := fmt.Sprintf("%s-%s", telemetryReportFilterRequest.StartDate, telemetryReportFilterRequest.EndDate)
	ginContext.Header("Content-Disposition", fmt.Sprintf("`attachment; filename=\"%s\"`", filename))
	ginContext.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile.Bytes())
}

func (telemetryHandler *Handler) IntervalReport(ginContext *gin.Context) {
	var telemetryIntervalFilterRequest model.TelemetryIntervalFilterRequest
	err := ginContext.ShouldBindBodyWithJSON(&telemetryIntervalFilterRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	telemetryResponses := telemetryHandler.telemetryService.IntervalReport(ginContext, &telemetryIntervalFilterRequest)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponse("Parameters has been fetched", telemetryResponses))
}
