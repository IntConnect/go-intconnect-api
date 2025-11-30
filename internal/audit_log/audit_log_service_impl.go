package audit_log

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	auditLogRepository Repository
	validatorService   validator.Service
	dbConnection       *gorm.DB
	viperConfig        *viper.Viper
	loggerInstance     *logrus.Logger
}

func NewService(auditLogRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, loggerInstance *logrus.Logger) *ServiceImpl {
	return &ServiceImpl{
		auditLogRepository: auditLogRepository,
		validatorService:   validatorService,
		dbConnection:       dbConnection,
		viperConfig:        viperConfig,
		loggerInstance:     loggerInstance,
	}
}

func (auditLogService *ServiceImpl) FindAll() []*model.AuditLogResponse {
	var allAuditLog []*model.AuditLogResponse
	err := auditLogService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		auditLogResponse, err := auditLogService.auditLogRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allAuditLog = helper.MapEntitiesIntoResponsesWithFunc[entity.AuditLog, *model.AuditLogResponse](auditLogResponse, mapper.FuncMapSimpleAuditable)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allAuditLog
}

func (auditLogService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.AuditLogResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var auditLogResponses []*model.AuditLogResponse
	var totalItems int64

	err := auditLogService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		auditLogEntities, total, err := auditLogService.auditLogRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		fmt.Println(auditLogEntities)
		auditLogResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.AuditLog, *model.AuditLogResponse](
			auditLogEntities,
			mapper.FuncMapSimpleAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Audit logs fetched successfully",
		auditLogResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat auditLog baru
func (auditLogService *ServiceImpl) Create(ginContext *gin.Context, createAuditLogRequest *model.CreateAuditLogRequest) *model.PaginatedResponse[*model.AuditLogResponse] {
	var paginationResp *model.PaginatedResponse[*model.AuditLogResponse]
	valErr := auditLogService.validatorService.ValidateStruct(createAuditLogRequest)
	auditLogService.validatorService.ParseValidationError(valErr, *createAuditLogRequest)
	err := auditLogService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		auditLogEntity := helper.MapCreateRequestIntoEntity[model.CreateAuditLogRequest, entity.AuditLog](createAuditLogRequest)
		auditLogEntity.SimpleAuditable = entity.NewSimpleAuditable(model.AUDIT_LOG_ACTOR_SYSTEM)
		err := auditLogService.auditLogRepository.Create(gormTransaction, auditLogEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = auditLogService.FindAllPagination(&paginationRequest)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}
