package parameter

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"math"
	"mime/multipart"

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

func (parameterService *ServiceImpl) FindAll() []*model.ParameterResponse {
	var parameterResponsesRequest []*model.ParameterResponse
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntities, err := parameterService.parameterRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		parameterResponsesRequest = helper.MapEntitiesIntoResponses[entity.Parameter, model.ParameterResponse](parameterEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return parameterResponsesRequest
}

func (parameterService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.ParameterResponse] {
	paginationResp := model.PaginationResponse[*model.ParameterResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allParameter []*model.ParameterResponse
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntities, totalItems, err := parameterService.parameterRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allParameter = helper.MapEntitiesIntoResponses[entity.Parameter, model.ParameterResponse](parameterEntities)
		paginationResp = model.PaginationResponse[*model.ParameterResponse]{
			Data:        allParameter,
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

// Create - Membuat parameter baru
func (parameterService *ServiceImpl) Create(ginContext *gin.Context, createParameterRequest *model.CreateParameterRequest, modelFile *multipart.FileHeader) {
	valErr := parameterService.validatorService.ValidateStruct(createParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *createParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parameterEntity := helper.MapCreateRequestIntoEntity[model.CreateParameterRequest, entity.Parameter](createParameterRequest)
		err := parameterService.parameterRepository.Create(gormTransaction, parameterEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
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

func (parameterService *ServiceImpl) Delete(ginContext *gin.Context, deleteParameterRequest *model.DeleteParameterRequest) {
	valErr := parameterService.validatorService.ValidateStruct(deleteParameterRequest)
	parameterService.validatorService.ParseValidationError(valErr, *deleteParameterRequest)
	err := parameterService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := parameterService.parameterRepository.Delete(gormTransaction, deleteParameterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
