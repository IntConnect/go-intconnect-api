package protocol_configuration

import (
	"fmt"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	protocolConfigurationRepository Repository
	validatorService                validator.Service
	dbConnection                    *gorm.DB
	viperConfig                     *viper.Viper
}

func NewService(protocolConfigurationRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		protocolConfigurationRepository: protocolConfigurationRepository,
		validatorService:                validatorService,
		dbConnection:                    dbConnection,
		viperConfig:                     viperConfig,
	}
}

func (protocolConfigurationService *ServiceImpl) FindAll() []*model.ProtocolConfigurationResponse {
	var allProtocolConfiguration []*model.ProtocolConfigurationResponse
	err := protocolConfigurationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		protocolConfigurationResponse, err := protocolConfigurationService.protocolConfigurationRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allProtocolConfiguration = mapper.MapProtocolConfigurationEntitiesIntoProtocolConfigurationResponses(protocolConfigurationResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allProtocolConfiguration
}

func (protocolConfigurationService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.ProtocolConfigurationResponse] {
	paginationResp := model.PaginationResponse[*model.ProtocolConfigurationResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allProtocolConfiguration []*model.ProtocolConfigurationResponse
	err := protocolConfigurationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		protocolConfigurationEntities, totalItems, err := protocolConfigurationService.protocolConfigurationRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allProtocolConfiguration = mapper.MapProtocolConfigurationEntitiesIntoProtocolConfigurationResponses(protocolConfigurationEntities)
		paginationResp = model.PaginationResponse[*model.ProtocolConfigurationResponse]{
			Data:        allProtocolConfiguration,
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

func (protocolConfigurationService *ServiceImpl) FindById(ginContext *gin.Context, protocolConfigurationId uint64) *model.ProtocolConfigurationResponse {
	var protocolConfigurationResponse *model.ProtocolConfigurationResponse
	err := protocolConfigurationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		protocolConfigurationEntity, err := protocolConfigurationService.protocolConfigurationRepository.FindById(gormTransaction, protocolConfigurationId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		protocolConfigurationResponse = mapper.MapProtocolConfigurationEntityIntoProtocolConfigurationResponse(protocolConfigurationEntity)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return protocolConfigurationResponse
}

// Create - Membuat protocolConfiguration baru
func (protocolConfigurationService *ServiceImpl) Create(ginContext *gin.Context, createProtocolConfigurationRequest *model.CreateProtocolConfigurationRequest) {
	valErr := protocolConfigurationService.validatorService.ValidateStruct(createProtocolConfigurationRequest)
	protocolConfigurationService.validatorService.ParseValidationError(valErr, *createProtocolConfigurationRequest)
	err := protocolConfigurationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		fmt.Println(createProtocolConfigurationRequest)

		protocolConfigurationEntity := mapper.MapCreateProtocolConfigurationRequestIntoProtocolConfigurationEntity(createProtocolConfigurationRequest)
		fmt.Println(protocolConfigurationEntity)
		err := protocolConfigurationService.protocolConfigurationRepository.Create(gormTransaction, protocolConfigurationEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (protocolConfigurationService *ServiceImpl) Update(ginContext *gin.Context, updateProtocolConfigurationRequest *model.UpdateProtocolConfigurationRequest) {
	valErr := protocolConfigurationService.validatorService.ValidateStruct(updateProtocolConfigurationRequest)
	protocolConfigurationService.validatorService.ParseValidationError(valErr, *updateProtocolConfigurationRequest)
	err := protocolConfigurationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		protocolConfiguration, err := protocolConfigurationService.protocolConfigurationRepository.FindById(gormTransaction, updateProtocolConfigurationRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateProtocolConfigurationRequestIntoProtocolConfigurationEntity(updateProtocolConfigurationRequest, protocolConfiguration)
		err = protocolConfigurationService.protocolConfigurationRepository.Update(gormTransaction, protocolConfiguration)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (protocolConfigurationService *ServiceImpl) Delete(ginContext *gin.Context, deleteProtocolConfigurationRequest *model.DeleteProtocolConfigurationRequest) {
	valErr := protocolConfigurationService.validatorService.ValidateStruct(deleteProtocolConfigurationRequest)
	protocolConfigurationService.validatorService.ParseValidationError(valErr, *deleteProtocolConfigurationRequest)
	err := protocolConfigurationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := protocolConfigurationService.protocolConfigurationRepository.Delete(gormTransaction, deleteProtocolConfigurationRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
