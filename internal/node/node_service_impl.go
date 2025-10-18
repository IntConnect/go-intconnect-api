package node

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

func (nodeService *ServiceImpl) FindAll() []*model.NodeResponse {
	var nodeResponsesRequest []*model.NodeResponse
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeEntities, err := nodeService.nodeRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		nodeResponsesRequest = mapper.MapNodeEntitiesIntoNodeResponses(nodeEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return nodeResponsesRequest
}

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
func (nodeService *ServiceImpl) Create(ginContext *gin.Context, createNodeRequest *model.CreateNodeRequest) {
	valErr := nodeService.validatorService.ValidateStruct(createNodeRequest)
	nodeService.validatorService.ParseValidationError(valErr, *createNodeRequest)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		nodeEntity := mapper.MapCreateNodeRequestIntoNodeEntity(createNodeRequest)
		err := nodeService.nodeRepository.Create(gormTransaction, nodeEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (nodeService *ServiceImpl) Update(ginContext *gin.Context, updateNodeRequest *model.UpdateNodeRequest) {
	valErr := nodeService.validatorService.ValidateStruct(updateNodeRequest)
	nodeService.validatorService.ParseValidationError(valErr, *updateNodeRequest)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		node, err := nodeService.nodeRepository.FindById(gormTransaction, updateNodeRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdateNodeRequestIntoNodeEntity(updateNodeRequest, node)
		err = nodeService.nodeRepository.Update(gormTransaction, node)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (nodeService *ServiceImpl) Delete(ginContext *gin.Context, deleteNodeRequest *model.DeleteNodeRequest) {
	valErr := nodeService.validatorService.ValidateStruct(deleteNodeRequest)
	nodeService.validatorService.ParseValidationError(valErr, *deleteNodeRequest)
	err := nodeService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := nodeService.nodeRepository.Delete(gormTransaction, deleteNodeRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
