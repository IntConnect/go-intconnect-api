package mqtt_broker

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	mqttBrokerRepository Repository
	auditLogService      auditLog.Service
	validatorService     validator.Service
	dbConnection         *gorm.DB
	viperConfig          *viper.Viper
}

func NewService(mqttBrokerRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		mqttBrokerRepository: mqttBrokerRepository,
		auditLogService:      auditLogService,
		validatorService:     validatorService,
		dbConnection:         dbConnection,
		viperConfig:          viperConfig,
	}
}

func (mqttBrokerService *ServiceImpl) FindAll() []*model.MqttBrokerResponse {
	var mqttBrokerResponsesRequest []*model.MqttBrokerResponse
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBrokerEntities, err := mqttBrokerService.mqttBrokerRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mqttBrokerResponsesRequest = helper.MapEntitiesIntoResponses[entity.MqttBroker, model.MqttBrokerResponse](mqttBrokerEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return mqttBrokerResponsesRequest
}

func (mqttBrokerService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MqttBrokerResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var mqttBrokerResponses []*model.MqttBrokerResponse
	var totalItems int64
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBrokerEntities, total, err := mqttBrokerService.mqttBrokerRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		mqttBrokerResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.MqttBroker, *model.MqttBrokerResponse](
			mqttBrokerEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"MQTT Broker fetched successfully",
		mqttBrokerResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat mqttBroker baru
func (mqttBrokerService *ServiceImpl) Create(ginContext *gin.Context, createMqttBrokerRequest *model.CreateMqttBrokerRequest) *model.PaginatedResponse[*model.MqttBrokerResponse] {
	jwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	ipAddress, _ := helper.ExtractRequestMeta(ginContext)
	var paginationResp *model.PaginatedResponse[*model.MqttBrokerResponse]
	valErr := mqttBrokerService.validatorService.ValidateStruct(createMqttBrokerRequest)
	mqttBrokerService.validatorService.ParseValidationError(valErr, *createMqttBrokerRequest)
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBrokerEntity := helper.MapCreateRequestIntoEntity[model.CreateMqttBrokerRequest, entity.MqttBroker](createMqttBrokerRequest)
		err := mqttBrokerService.mqttBrokerRepository.Create(gormTransaction, mqttBrokerEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = mqttBrokerService.FindAllPagination(&paginationRequest)
		mqttBrokerService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      jwtClaims.Id,
			Action:      model.AUDIT_LOG_CREATE,
			Feature:     model.AUDIT_LOG_FEATURE_MQTT_BROKER,
			Description: "",
			Before:      nil,
			After:       mqttBrokerEntity,
			IpAddress:   ipAddress,
		})
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

func (mqttBrokerService *ServiceImpl) Update(ginContext *gin.Context, updateMqttBrokerRequest *model.UpdateMqttBrokerRequest) *model.PaginatedResponse[*model.MqttBrokerResponse] {
	var paginationResp *model.PaginatedResponse[*model.MqttBrokerResponse]
	valErr := mqttBrokerService.validatorService.ValidateStruct(updateMqttBrokerRequest)
	mqttBrokerService.validatorService.ParseValidationError(valErr, *updateMqttBrokerRequest)
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBroker, err := mqttBrokerService.mqttBrokerRepository.FindById(gormTransaction, updateMqttBrokerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateMqttBrokerRequest, mqttBroker)
		err = mqttBrokerService.mqttBrokerRepository.Update(gormTransaction, mqttBroker)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = mqttBrokerService.FindAllPagination(&paginationRequest)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

func (mqttBrokerService *ServiceImpl) Delete(ginContext *gin.Context, deleteMqttBrokerRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.MqttBrokerResponse] {
	var paginationResp *model.PaginatedResponse[*model.MqttBrokerResponse]

	valErr := mqttBrokerService.validatorService.ValidateStruct(deleteMqttBrokerRequest)
	mqttBrokerService.validatorService.ParseValidationError(valErr, *deleteMqttBrokerRequest)
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := mqttBrokerService.mqttBrokerRepository.Delete(gormTransaction, deleteMqttBrokerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = mqttBrokerService.FindAllPagination(&paginationRequest)

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}
