package parameter

import (
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
	parameterRepository Repository
	validatorService    validator.Service
	dbConnection        *gorm.DB
	viperConfig         *viper.Viper
}

func NewService(parameterRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		parameterRepository: parameterRepository,
		validatorService:    validatorService,
		dbConnection:        dbConnection,
		viperConfig:         viperConfig,
	}
}

func (parameterService *ServiceImpl) FindAll(parameterFilterRequest *model.ParameterFilterRequest) []*model.ParameterResponse {
	var parameterResponsesRequest []*model.ParameterResponse
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntities, err := parameterService.parameterRepository.FindAll(gormTransaction, parameterFilterRequest)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		parameterResponsesRequest = helper.MapEntitiesIntoResponsesWithFunc[
			*entity.Parameter,
			*model.ParameterResponse,
		](
			parameterEntities,
			mapper.FuncMapAuditable[*entity.Parameter, *model.ParameterResponse],
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return parameterResponsesRequest
}

func (parameterService *ServiceImpl) FindDependencyParameter() *model.ParameterDependency {
	var parameterDependency = &model.ParameterDependency{}
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var machineEntities []entity.Machine
		var mqttTopicEntities []entity.MqttTopic
		err := gormTransaction.Find(&machineEntities).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Find(&mqttTopicEntities).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		parameterDependency.MachineResponses = helper.MapEntitiesIntoResponses[entity.Machine, model.MachineResponse](machineEntities)
		parameterDependency.MqttTopicResponses = helper.MapEntitiesIntoResponses[entity.MqttTopic, model.MqttTopicResponse](mqttTopicEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return parameterDependency
}

func (parameterService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.ParameterResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var parameterResponses []*model.ParameterResponse
	var totalItems int64
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntities, total, err := parameterService.parameterRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery)
		parameterResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.Parameter, *model.ParameterResponse](
			parameterEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Parameters fetched successfully",
		parameterResponses,
		paginationReq,
		totalItems,
	)
}

func (parameterService *ServiceImpl) FindById(ginContext *gin.Context, parameterId uint64) *model.ParameterResponse {
	var parameterResponse *model.ParameterResponse
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntity, err := parameterService.parameterRepository.FindById(gormTransaction, parameterId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		parameterResponse = helper.MapEntityIntoResponse[
			*entity.Parameter,
			*model.ParameterResponse,
		](
			parameterEntity,
			mapper.FuncMapAuditable[*entity.Parameter, *model.ParameterResponse],
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return parameterResponse
}

// Create - Membuat parameter baru
func (parameterService *ServiceImpl) Create(ginContext *gin.Context, createParameterRequest *model.CreateParameterRequest) *model.PaginatedResponse[*model.ParameterResponse] {
	var parameterResponses *model.PaginatedResponse[*model.ParameterResponse]
	valErr := parameterService.validatorService.ValidateStruct(createParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *createParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntity := helper.MapCreateRequestIntoEntity[model.CreateParameterRequest, entity.Parameter](createParameterRequest)
		err := parameterService.parameterRepository.Create(gormTransaction, parameterEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		parameterResponses = parameterService.FindAllPagination(&paginationRequest)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return parameterResponses
}

func (parameterService *ServiceImpl) Update(ginContext *gin.Context, updateParameterRequest *model.UpdateParameterRequest) {
	valErr := parameterService.validatorService.ValidateStruct(updateParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *updateParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameter, err := parameterService.parameterRepository.FindById(gormTransaction, updateParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateParameterRequest, parameter)
		err = parameterService.parameterRepository.Update(gormTransaction, parameter)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (parameterService *ServiceImpl) UpdateOperation(ginContext *gin.Context, updateParameterRequest *model.ManageParameterOperationRequest) *model.PaginatedResponse[*model.ParameterResponse] {
	valErr := parameterService.validatorService.ValidateStruct(updateParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *updateParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntity, err := parameterService.parameterRepository.FindById(gormTransaction, updateParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		var parameterOperationEntities []*entity.ParameterOperation
		for _, parameterOperation := range updateParameterRequest.ParameterOperations {
			parameterOperationEntity := helper.MapCreateRequestIntoEntity[model.ParameterOperationRequest, entity.ParameterOperation](parameterOperation)
			parameterOperationEntities = append(parameterOperationEntities, parameterOperationEntity)
		}
		parameterEntity.ParameterOperation = parameterOperationEntities
		err = parameterService.parameterRepository.Update(gormTransaction, parameterEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	var paginatedResp *model.PaginatedResponse[*model.ParameterResponse]
	paginationRequest := model.NewPaginationRequest()
	paginatedResp = parameterService.FindAllPagination(&paginationRequest)
	return paginatedResp
}

func (parameterService *ServiceImpl) Delete(ginContext *gin.Context, deleteParameterRequest *model.DeleteResourceGeneralRequest) {
	valErr := parameterService.validatorService.ValidateStruct(deleteParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *deleteParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := parameterService.parameterRepository.Delete(gormTransaction, deleteParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
