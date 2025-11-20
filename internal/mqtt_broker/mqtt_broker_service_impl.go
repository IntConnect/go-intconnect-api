package mqtt_broker

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	mqttBrokerRepository Repository
	validatorService     validator.Service
	dbConnection         *gorm.DB
	viperConfig          *viper.Viper
}

func NewService(mqttBrokerRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		mqttBrokerRepository: mqttBrokerRepository,
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

func (mqttBrokerService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.MqttBrokerResponse] {
	paginationResp := model.PaginationResponse[*model.MqttBrokerResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allMqttBroker []*model.MqttBrokerResponse
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBrokerEntities, totalItems, err := mqttBrokerService.mqttBrokerRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allMqttBroker = helper.MapEntitiesIntoResponses[entity.MqttBroker, model.MqttBrokerResponse](mqttBrokerEntities)
		paginationResp = model.PaginationResponse[*model.MqttBrokerResponse]{
			Data:        allMqttBroker,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: paginationReq.Page,
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}

// Create - Membuat mqttBroker baru
func (mqttBrokerService *ServiceImpl) Create(ginContext *gin.Context, createMqttBrokerRequest *model.CreateMqttBrokerRequest) {
	valErr := mqttBrokerService.validatorService.ValidateStruct(createMqttBrokerRequest)
	mqttBrokerService.validatorService.ParseValidationError(valErr, *createMqttBrokerRequest)
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBrokerEntity := helper.MapCreateRequestIntoEntity[model.CreateMqttBrokerRequest, entity.MqttBroker](createMqttBrokerRequest)
		err := mqttBrokerService.mqttBrokerRepository.Create(gormTransaction, mqttBrokerEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (mqttBrokerService *ServiceImpl) Update(ginContext *gin.Context, updateMqttBrokerRequest *model.UpdateMqttBrokerRequest) {
	valErr := mqttBrokerService.validatorService.ValidateStruct(updateMqttBrokerRequest)
	mqttBrokerService.validatorService.ParseValidationError(valErr, *updateMqttBrokerRequest)
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		mqttBroker, err := mqttBrokerService.mqttBrokerRepository.FindById(gormTransaction, updateMqttBrokerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateMqttBrokerRequest, mqttBroker)
		err = mqttBrokerService.mqttBrokerRepository.Update(gormTransaction, mqttBroker)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (mqttBrokerService *ServiceImpl) Delete(ginContext *gin.Context, deleteMqttBrokerRequest *model.DeleteResourceGeneralRequest) {
	valErr := mqttBrokerService.validatorService.ValidateStruct(deleteMqttBrokerRequest)
	mqttBrokerService.validatorService.ParseValidationError(valErr, *deleteMqttBrokerRequest)
	err := mqttBrokerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := mqttBrokerService.mqttBrokerRepository.Delete(gormTransaction, deleteMqttBrokerRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
