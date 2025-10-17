package pipeline

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/node"
	pipelineEdge "go-intconnect-api/internal/pipeline_edge"
	pipelineNode "go-intconnect-api/internal/pipeline_node"
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
	pipelineRepository     Repository
	pipelineNodeRepository pipelineNode.Repository
	pipelineEdgeRepository pipelineEdge.Repository
	nodeRepository         node.Repository
	validatorService       validator.Service
	dbConnection           *gorm.DB
	viperConfig            *viper.Viper
}

func NewService(pipelineRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	pipelineNodeRepository pipelineNode.Repository,
	pipelineEdgeRepository pipelineEdge.Repository,
	nodeRepository node.Repository,

) *ServiceImpl {
	return &ServiceImpl{
		pipelineRepository:     pipelineRepository,
		validatorService:       validatorService,
		dbConnection:           dbConnection,
		viperConfig:            viperConfig,
		pipelineNodeRepository: pipelineNodeRepository,
		pipelineEdgeRepository: pipelineEdgeRepository,
		nodeRepository:         nodeRepository,
	}
}

func (pipelineService *ServiceImpl) FindAll() []*model.PipelineResponse {
	var pipelineResponsesRequest []*model.PipelineResponse
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntities, err := pipelineService.pipelineRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		pipelineResponsesRequest = mapper.MapPipelineEntitiesIntoPipelineResponses(pipelineEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return pipelineResponsesRequest
}

func (pipelineService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.PipelineResponse] {
	paginationResp := model.PaginationResponse[*model.PipelineResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allPipeline []*model.PipelineResponse
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntities, totalItems, err := pipelineService.pipelineRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allPipeline = mapper.MapPipelineEntitiesIntoPipelineResponses(pipelineEntities)
		paginationResp = model.PaginationResponse[*model.PipelineResponse]{
			Data:        allPipeline,
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

// Create - Membuat pipeline baru
func (pipelineService *ServiceImpl) Create(ginContext *gin.Context, createPipelineRequest *model.CreatePipelineRequest) {
	valErr := pipelineService.validatorService.ValidateStruct(createPipelineRequest)
	pipelineService.validatorService.ParseValidationError(valErr, *createPipelineRequest)

	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntity := mapper.MapCreatePipelineRequestIntoPipelineEntity(createPipelineRequest)
		pipelineEntity.Auditable = entity.NewAuditable("System")
		err := pipelineService.pipelineRepository.Create(gormTransaction, pipelineEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		mapOfTemporaryIdWithDatabaseId := map[string]uint64{}
		var nodeIds []uint64

		for i := range pipelineEntity.PipelineNode {
			nodeIds = append(nodeIds, pipelineEntity.PipelineNode[i].NodeId)
		}

		nodeEntities, err := pipelineService.nodeRepository.FindBatchById(gormTransaction, nodeIds)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		if len(nodeEntities) != len(nodeIds) {
			return fmt.Errorf("not all nodes found: expected %d, got %d", len(nodeIds), len(nodeEntities))
		}

		// ✅ Fixed: Use index to work with pointers
		for i := range pipelineEntity.PipelineNode {
			pipelineNodeEntity := pipelineEntity.PipelineNode[i]
			pipelineNodeEntity.PipelineId = pipelineEntity.Id
			pipelineNodeEntity.Auditable = entity.NewAuditable("System")

			err := pipelineService.pipelineNodeRepository.Create(gormTransaction, pipelineNodeEntity)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))

			mapOfTemporaryIdWithDatabaseId[pipelineNodeEntity.TempId] = pipelineNodeEntity.Id
		}

		// ✅ Fixed: Same for edges
		for i := range pipelineEntity.PipelineEdge {
			pipelineEdgeEntity := pipelineEntity.PipelineEdge[i]
			pipelineEdgeEntity.PipelineId = pipelineEntity.Id
			pipelineEdgeEntity.SourceNodeId = mapOfTemporaryIdWithDatabaseId[pipelineEdgeEntity.SourceNodeTempId]
			pipelineEdgeEntity.TargetNodeId = mapOfTemporaryIdWithDatabaseId[pipelineEdgeEntity.TargetNodeTempId]

			err := pipelineService.pipelineEdgeRepository.Create(gormTransaction, pipelineEdgeEntity)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}

		return nil
	})

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
func (pipelineService *ServiceImpl) Update(ginContext *gin.Context, updatePipelineRequest *model.UpdatePipelineRequest) {
	valErr := pipelineService.validatorService.ValidateStruct(updatePipelineRequest)
	pipelineService.validatorService.ParseValidationError(valErr, *updatePipelineRequest)
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipeline, err := pipelineService.pipelineRepository.FindById(gormTransaction, updatePipelineRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapUpdatePipelineRequestIntoPipelineEntity(updatePipelineRequest, pipeline)
		err = pipelineService.pipelineRepository.Update(gormTransaction, pipeline)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (pipelineService *ServiceImpl) Delete(ginContext *gin.Context, deletePipelineRequest *model.DeletePipelineRequest) {
	valErr := pipelineService.validatorService.ValidateStruct(deletePipelineRequest)
	pipelineService.validatorService.ParseValidationError(valErr, *deletePipelineRequest)
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := pipelineService.pipelineRepository.Delete(gormTransaction, deletePipelineRequest.ID)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
