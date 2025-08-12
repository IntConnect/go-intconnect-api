package node

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"gorm.io/gorm"
	"math"
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

// Create - Membuat node baru
func (nodeService *ServiceImpl) FindAll() []*model.NodeResponse {
	var allNode []*model.NodeResponse
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeResponse, err := nodeService.nodeRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		allNode = mapper.MapNodeEntitiesIntoNodeResponses(nodeResponse)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return allNode
}

// Create - Membuat node baru
func (nodeService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.NodeResponse] {
	paginationResp := model.PaginationResponse[*model.NodeResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allNode []*model.NodeResponse
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeEntities, totalItems, err := nodeService.nodeRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allNode = mapper.MapNodeEntitiesIntoNodeResponses(nodeEntities)
		paginationResp = model.PaginationResponse[*model.NodeResponse]{
			Data:        allNode,
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
func (nodeService *ServiceImpl) Create(ginContext *gin.Context, createNodeDto *model.CreateNodeDto) {
	valErr := nodeService.validatorService.ValidateStruct(createNodeDto)
	nodeService.validatorService.ParseValidationError(valErr, *createNodeDto)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeEntity := mapper.MapCreateNodeDtoIntoNodeEntity(createNodeDto)
		err := nodeService.nodeRepository.Create(gormTransaction, nodeEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (nodeService *ServiceImpl) Update(ginContext *gin.Context, updateNodeDto *model.UpdateNodeDto) {
	valErr := nodeService.validatorService.ValidateStruct(updateNodeDto)
	nodeService.validatorService.ParseValidationError(valErr, *updateNodeDto)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		node, err := nodeService.nodeRepository.FindById(gormTransaction, updateNodeDto.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateNodeDtoIntoNodeEntity(updateNodeDto, node)
		err = nodeService.nodeRepository.Update(gormTransaction, node)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (nodeService *ServiceImpl) Delete(ginContext *gin.Context, deleteNodeDto *model.DeleteNodeDto) {
	valErr := nodeService.validatorService.ValidateStruct(deleteNodeDto)
	nodeService.validatorService.ParseValidationError(valErr, *deleteNodeDto)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := nodeService.nodeRepository.Delete(gormTransaction, deleteNodeDto.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
