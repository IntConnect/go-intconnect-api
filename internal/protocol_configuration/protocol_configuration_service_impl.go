package protocol_configuration

import (
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
	nodeRepository   Repository
	validatorService validator.Service
	dbConnection     *gorm.DB
	viperConfig      *viper.Viper
}

func NewService(nodeRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		nodeRepository:   nodeRepository,
		validatorService: validatorService,
		dbConnection:     dbConnection,
		viperConfig:      viperConfig,
	}
}

func (nodeService *ServiceImpl) FindAll() []*model.ProtocolConfigurationResponse {
	var allProtocolConfiguration []*model.ProtocolConfigurationResponse
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeResponse, err := nodeService.nodeRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allProtocolConfiguration = mapper.MapProtocolConfigurationEntitiesIntoProtocolConfigurationResponses(nodeResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allProtocolConfiguration
}

func (nodeService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.ProtocolConfigurationResponse] {
	paginationResp := model.PaginationResponse[*model.ProtocolConfigurationResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allProtocolConfiguration []*model.ProtocolConfigurationResponse
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeEntities, totalItems, err := nodeService.nodeRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allProtocolConfiguration = mapper.MapProtocolConfigurationEntitiesIntoProtocolConfigurationResponses(nodeEntities)
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

// Create - Membuat node baru
func (nodeService *ServiceImpl) Create(ginContext *gin.Context, createProtocolConfigurationDto *model.CreateProtocolConfigurationDto) {
	valErr := nodeService.validatorService.ValidateStruct(createProtocolConfigurationDto)
	nodeService.validatorService.ParseValidationError(valErr, *createProtocolConfigurationDto)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeEntity := mapper.MapCreateProtocolConfigurationDtoIntoProtocolConfigurationEntity(createProtocolConfigurationDto)
		err := nodeService.nodeRepository.Create(gormTransaction, nodeEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (nodeService *ServiceImpl) Update(ginContext *gin.Context, updateProtocolConfigurationDto *model.UpdateProtocolConfigurationDto) {
	valErr := nodeService.validatorService.ValidateStruct(updateProtocolConfigurationDto)
	nodeService.validatorService.ParseValidationError(valErr, *updateProtocolConfigurationDto)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		node, err := nodeService.nodeRepository.FindById(gormTransaction, updateProtocolConfigurationDto.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateProtocolConfigurationDtoIntoProtocolConfigurationEntity(updateProtocolConfigurationDto, node)
		err = nodeService.nodeRepository.Update(gormTransaction, node)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (nodeService *ServiceImpl) Delete(ginContext *gin.Context, deleteProtocolConfigurationDto *model.DeleteProtocolConfigurationDto) {
	valErr := nodeService.validatorService.ValidateStruct(deleteProtocolConfigurationDto)
	nodeService.validatorService.ParseValidationError(valErr, *deleteProtocolConfigurationDto)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := nodeService.nodeRepository.Delete(gormTransaction, deleteProtocolConfigurationDto.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
