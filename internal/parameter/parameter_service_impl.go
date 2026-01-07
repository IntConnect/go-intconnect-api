package parameter

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	parameterOperation "go-intconnect-api/internal/parameter_operation"
	processedParameterSequence "go-intconnect-api/internal/processed_parameter_sequence"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	parameterRepository                  Repository
	parameterOperationRepository         parameterOperation.Repository
	processedParameterSequenceRepository processedParameterSequence.Repository
	auditLogService                      auditLog.Service
	validatorService                     validator.Service
	dbConnection                         *gorm.DB
	viperConfig                          *viper.Viper
}

func NewService(parameterRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	auditLogService auditLog.Service,
	parameterOperationRepository parameterOperation.Repository,
	processedParameterSequenceRepository processedParameterSequence.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		parameterRepository:                  parameterRepository,
		validatorService:                     validatorService,
		dbConnection:                         dbConnection,
		viperConfig:                          viperConfig,
		auditLogService:                      auditLogService,
		parameterOperationRepository:         parameterOperationRepository,
		processedParameterSequenceRepository: processedParameterSequenceRepository,
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
			mapper.FuncMapParameter,
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return parameterResponse
}

// Create - Membuat parameter baru
func (parameterService *ServiceImpl) Create(ginContext *gin.Context, createParameterRequest *model.CreateParameterRequest) {
	valErr := parameterService.validatorService.ValidateStruct(createParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *createParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntity := helper.MapCreateRequestIntoEntity[model.CreateParameterRequest, entity.Parameter](createParameterRequest)
		err := parameterService.parameterRepository.Create(gormTransaction, parameterEntity)
		var processedParameterSequences []*entity.ProcessedParameterSequence
		for _, processedParameterSequenceRequest := range createParameterRequest.ProcessedParameterSequence {
			processedParameterSequences = append(processedParameterSequences, &entity.ProcessedParameterSequence{
				ParentParameterId: parameterEntity.Id,
				ParameterId:       processedParameterSequenceRequest.ParameterId,
				Sequence:          processedParameterSequenceRequest.Sequence,
				Type:              processedParameterSequenceRequest.Type,
			})
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = parameterService.processedParameterSequenceRepository.CreateBatch(gormTransaction, processedParameterSequences)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		auditPayload := parameterService.auditLogService.Build(
			nil,             // before entity
			parameterEntity, // after entity
			nil,
			"",
		)
		err = parameterService.auditLogService.Record(ginContext,
			model.AUDIT_LOG_CREATE,
			model.AUDIT_LOG_FEATURE_PARAMETER,
			auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (parameterService *ServiceImpl) Update(ginContext *gin.Context, updateParameterRequest *model.UpdateParameterRequest) {
	valErr := parameterService.validatorService.ValidateStruct(updateParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *updateParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {

		parameterEntity, err := parameterService.parameterRepository.FindById(gormTransaction, updateParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity(updateParameterRequest, parameterEntity)
		beforeParameterEntity := *parameterEntity
		if !parameterEntity.IsAutomatic {
			parameterEntity.MqttTopicId = nil
		}
		err = parameterService.parameterRepository.Update(gormTransaction, parameterEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := parameterService.auditLogService.Build(
			beforeParameterEntity, // before entity
			parameterEntity,       // after entity
			nil,
			"",
		)
		err = parameterService.auditLogService.Record(ginContext,
			model.AUDIT_LOG_UPDATE,
			model.AUDIT_LOG_FEATURE_PARAMETER,
			auditPayload)
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
		deletedParameterOperationEntities, err := parameterService.parameterOperationRepository.FindBatchById(gormTransaction, updateParameterRequest.Deleted)
		if len(deletedParameterOperationEntities) != len(updateParameterRequest.Deleted) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusNotFound, exception.ErrSomeResourceNotFound))
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		var createdParameterOperationEntities []*entity.ParameterOperation
		for _, parameterOperationRequest := range updateParameterRequest.Created {
			parameterOperationEntity := helper.MapCreateRequestIntoEntity[model.ParameterOperationRequest, entity.ParameterOperation](parameterOperationRequest)
			parameterOperationEntity.ParameterId = parameterEntity.Id
			createdParameterOperationEntities = append(createdParameterOperationEntities, parameterOperationEntity)
		}
		var updateParameterOperationEntities []*entity.ParameterOperation
		for _, parameterOperationRequest := range updateParameterRequest.Updated {
			parameterOperationEntity := helper.MapCreateRequestIntoEntity[model.ParameterOperationRequest, entity.ParameterOperation](parameterOperationRequest)
			parameterOperationEntity.ParameterId = parameterEntity.Id
			updateParameterOperationEntities = append(updateParameterOperationEntities, parameterOperationEntity)
		}
		if len(updateParameterRequest.Deleted) > 0 {
			err = parameterService.parameterOperationRepository.DeleteBatchById(gormTransaction, updateParameterRequest.Deleted)
		}
		if len(updateParameterOperationEntities) > 0 {
			for _, operationParameterEntity := range updateParameterOperationEntities {
				err = parameterService.parameterOperationRepository.Update(gormTransaction, operationParameterEntity)
				helper.CheckErrorOperation(err, exception.ParseGormError(err))
			}
		}
		if len(createdParameterOperationEntities) > 0 {
			err = parameterService.parameterOperationRepository.CreateBatch(gormTransaction, createdParameterOperationEntities)

		}
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
		parameterEntity, err := parameterService.parameterRepository.FindById(gormTransaction, deleteParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = parameterService.parameterRepository.Delete(gormTransaction, deleteParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := parameterService.auditLogService.Build(
			parameterEntity,
			nil,
			nil,
			deleteParameterRequest.Reason,
		)
		err = parameterService.auditLogService.Record(ginContext,
			model.AUDIT_LOG_DELETE,
			model.AUDIT_LOG_FEATURE_PARAMETER,
			auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
