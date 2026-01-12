package alarm_log

import (
	"go-intconnect-api/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	FindAll() []*model.AlarmLogResponse
	FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.AlarmLogResponse]
	Update(ginContext *gin.Context, updateAlarmLogRequest *model.UpdateAlarmLogRequest) *model.PaginatedResponse[*model.AlarmLogResponse]
	FindByMachineId(ginContext *gin.Context, machineId uint64) []*model.AlarmLogResponse
}
