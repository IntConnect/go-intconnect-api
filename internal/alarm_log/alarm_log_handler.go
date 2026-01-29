package alarm_log

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	alarmLogService Service
	viperConfig     *viper.Viper
}

func NewHandler(alarmLogService Service, viperConfig *viper.Viper) *Handler {
	return &Handler{
		alarmLogService: alarmLogService,
		viperConfig:     viperConfig,
	}
}

func (alarmLogHandler *Handler) FindAllAlarmLog(ginContext *gin.Context) {
	alarmLogResponses := alarmLogHandler.alarmLogService.FindAll()
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("AlarmLog has been fetched", alarmLogResponses))
}

func (alarmLogHandler *Handler) FindAllAlarmLogPagination(ginContext *gin.Context) {
	var paginationReq model.PaginationRequest
	err := ginContext.ShouldBindQuery(&paginationReq)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	paginatedResponse := alarmLogHandler.alarmLogService.FindAllPagination(&paginationReq)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (alarmLogHandler *Handler) UpdateAlarmLog(ginContext *gin.Context) {
	var updateAlarmLogRequest model.UpdateAlarmLogRequest
	err := ginContext.ShouldBindBodyWithJSON(&updateAlarmLogRequest)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	alarmLogId := ginContext.Param("id")
	parsedAlarmLogId, err := strconv.ParseUint(alarmLogId, 10, 64)
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	updateAlarmLogRequest.Id = parsedAlarmLogId
	paginatedResponse := alarmLogHandler.alarmLogService.Update(ginContext, &updateAlarmLogRequest)
	ginContext.JSON(http.StatusOK, paginatedResponse)
}

func (alarmLogHandler *Handler) FindAlarmLogByMachineId(ginContext *gin.Context) {
	machineId := ginContext.Param("id")
	parsedMachineId, err := strconv.ParseUint(machineId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	alarmLogResponses := alarmLogHandler.alarmLogService.FindByMachineId(ginContext, parsedMachineId, false)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Alarm log fetch successfully", alarmLogResponses))
}

func (alarmLogHandler *Handler) FindMinimalAlarmLogByMachineId(ginContext *gin.Context) {
	machineId := ginContext.Param("id")
	parsedMachineId, err := strconv.ParseUint(machineId, 10, 64)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	alarmLogResponses := alarmLogHandler.alarmLogService.FindByMachineId(ginContext, parsedMachineId, true)
	ginContext.JSON(http.StatusOK, helper.NewSuccessResponseWithEntries("Alarm log fetch successfully", alarmLogResponses))
}
