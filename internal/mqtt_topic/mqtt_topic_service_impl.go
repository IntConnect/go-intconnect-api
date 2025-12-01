package mqtt_topic

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/machine"
	"go-intconnect-api/internal/model"
	mqttBroker "go-intconnect-api/internal/mqtt_broker"
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
	mqttTopicRepository Repository
	auditLogService     auditLog.Service
	validatorService    validator.Service
	dbConnection        *gorm.DB
	viperConfig         *viper.Viper
	loggerInstance      *logrus.Logger
}

func NewService(mqttTopicRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, loggerInstance *logrus.Logger,
	auditLogService auditLog.Service,
	machineRepository machine.Repository, mqttBrokerRepository mqttBroker.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		mqttTopicRepository: mqttTopicRepository,
		validatorService:    validatorService,
		dbConnection:        dbConnection,
		viperConfig:         viperConfig,
		loggerInstance:      loggerInstance,
		auditLogService:     auditLogService,
	}
}

// Create - Membuat mqttTopic baru
func (mqttTopicService *ServiceImpl) FindAll() []*model.MqttTopicResponse {
	var allMqttTopic []*model.MqttTopicResponse
	err := mqttTopicService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttTopicResponse, err := mqttTopicService.mqttTopicRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allMqttTopic = helper.MapEntitiesIntoResponses[entity.MqttTopic, *model.MqttTopicResponse](mqttTopicResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allMqttTopic
}

func (mqttTopicService *ServiceImpl) FindDependency() *model.MqttTopicDependency {
	var mqttTopicDependency = &model.MqttTopicDependency{}
	err := mqttTopicService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var machineEntities []entity.Machine
		var mqttBrokerEntities []entity.MqttBroker
		err := gormTransaction.Find(&machineEntities).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Find(&mqttBrokerEntities).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mqttTopicDependency.MachineResponses = helper.MapEntitiesIntoResponses[entity.Machine, model.MachineResponse](machineEntities)
		mqttTopicDependency.MqttBrokerResponses = helper.MapEntitiesIntoResponses[entity.MqttBroker, model.MqttBrokerResponse](mqttBrokerEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return mqttTopicDependency
}

// Create - Membuat mqttTopic baru
func (mqttTopicService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.MqttTopicResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var mqttTopicResponses []*model.MqttTopicResponse
	var totalItems int64

	err := mqttTopicService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttTopicEntities, total, err := mqttTopicService.mqttTopicRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		mqttTopicResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.MqttTopic, *model.MqttTopicResponse](
			mqttTopicEntities,
			mapper.FuncMapAuditable,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Mqtt topic fetched successfully",
		mqttTopicResponses,
		paginationReq,
		totalItems,
	)
}

// Create - Membuat mqttTopic baru
func (mqttTopicService *ServiceImpl) Create(ginContext *gin.Context, createMqttTopicRequest *model.CreateMqttTopicRequest) *model.PaginatedResponse[*model.MqttTopicResponse] {
	var paginatedResp *model.PaginatedResponse[*model.MqttTopicResponse]
	valErr := mqttTopicService.validatorService.ValidateStruct(createMqttTopicRequest)
	mqttTopicService.validatorService.ParseValidationError(valErr, *createMqttTopicRequest)
	err := mqttTopicService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttTopicEntity := helper.MapCreateRequestIntoEntity[model.CreateMqttTopicRequest, entity.MqttTopic](createMqttTopicRequest)
		mqttTopicEntity.Auditable = entity.NewAuditable("Administrator")
		err := mqttTopicService.mqttTopicRepository.Create(gormTransaction, mqttTopicEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = mqttTopicService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (mqttTopicService *ServiceImpl) Update(ginContext *gin.Context, updateMqttTopicRequest *model.UpdateMqttTopicRequest) *model.PaginatedResponse[*model.MqttTopicResponse] {
	var paginatedResp *model.PaginatedResponse[*model.MqttTopicResponse]
	valErr := mqttTopicService.validatorService.ValidateStruct(updateMqttTopicRequest)
	mqttTopicService.validatorService.ParseValidationError(valErr, *updateMqttTopicRequest)
	err := mqttTopicService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttTopic, err := mqttTopicService.mqttTopicRepository.FindById(gormTransaction, updateMqttTopicRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateMqttTopicRequest, mqttTopic)
		err = mqttTopicService.mqttTopicRepository.Update(gormTransaction, mqttTopic)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = mqttTopicService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (mqttTopicService *ServiceImpl) Delete(ginContext *gin.Context, deleteMqttTopicRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.MqttTopicResponse] {
	var paginatedResp *model.PaginatedResponse[*model.MqttTopicResponse]
	valErr := mqttTopicService.validatorService.ValidateStruct(deleteMqttTopicRequest)
	mqttTopicService.validatorService.ParseValidationError(valErr, *deleteMqttTopicRequest)
	err := mqttTopicService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := mqttTopicService.mqttTopicRepository.Delete(gormTransaction, deleteMqttTopicRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = mqttTopicService.FindAllPagination(&paginationRequest)
	return paginatedResp
}
