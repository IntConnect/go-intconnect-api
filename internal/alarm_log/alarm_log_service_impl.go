package alarm_log

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	alarmLogRepository Repository
	validatorService   validator.Service
	dbConnection       *gorm.DB
	viperConfig        *viper.Viper
}

func NewService(alarmLogRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		alarmLogRepository: alarmLogRepository,
		validatorService:   validatorService,
		dbConnection:       dbConnection,
		viperConfig:        viperConfig,
	}
}

func (alarmLogService *ServiceImpl) FindAll() []*model.AlarmLogResponse {
	var allAlarmLog []*model.AlarmLogResponse
	err := alarmLogService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		alarmLogResponse, err := alarmLogService.alarmLogRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allAlarmLog = helper.MapEntitiesIntoResponsesWithFunc[*entity.AlarmLog, *model.AlarmLogResponse](alarmLogResponse, mapper.FuncMapAlarmLog)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allAlarmLog
}

func (alarmLogService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.AlarmLogResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var alarmLogResponses []*model.AlarmLogResponse
	var totalItems int64

	err := alarmLogService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		alarmLogEntities, total, err := alarmLogService.alarmLogRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		alarmLogResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.AlarmLog, *model.AlarmLogResponse](alarmLogEntities, mapper.FuncMapAlarmLog)

		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Alarm logs fetched successfully",
		alarmLogResponses,
		paginationReq,
		totalItems,
	)
}

func (alarmLogService *ServiceImpl) Update(ginContext *gin.Context, updateAlarmLogRequest *model.UpdateAlarmLogRequest) *model.PaginatedResponse[*model.AlarmLogResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginatedResp *model.PaginatedResponse[*model.AlarmLogResponse]
	valErr := alarmLogService.validatorService.ValidateStruct(updateAlarmLogRequest)
	alarmLogService.validatorService.ParseValidationError(valErr, *updateAlarmLogRequest)
	err := alarmLogService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		alarmLogEntity, err := alarmLogService.alarmLogRepository.FindById(gormTransaction, updateAlarmLogRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		alarmLogEntity.Note = updateAlarmLogRequest.Note
		alarmLogEntity.AcknowledgedAt = helper.TakePointer(time.Now())
		alarmLogEntity.AcknowledgedBy = helper.TakePointer(userJwtClaims.Id)
		alarmLogEntity.Status = "Acknowledged"
		err = alarmLogService.alarmLogRepository.Update(gormTransaction, alarmLogEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = alarmLogService.FindAllPagination(&paginationRequest)
	return paginatedResp
}
