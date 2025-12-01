package pipeline

import (
	"encoding/json"
	"fmt"
	"go-intconnect-api/configs"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/node"
	pipelineEdge "go-intconnect-api/internal/pipeline_edge"
	pipelineNode "go-intconnect-api/internal/pipeline_node"
	protocolConfiguration "go-intconnect-api/internal/protocol_configuration"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	pipelineRepository              Repository
	pipelineNodeRepository          pipelineNode.Repository
	pipelineEdgeRepository          pipelineEdge.Repository
	nodeRepository                  node.Repository
	validatorService                validator.Service
	dbConnection                    *gorm.DB
	viperConfig                     *viper.Viper
	redisInstance                   *configs.RedisInstance
	protocolConfigurationRepository protocolConfiguration.Repository
}

func NewService(pipelineRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper,
	pipelineNodeRepository pipelineNode.Repository,
	pipelineEdgeRepository pipelineEdge.Repository,
	nodeRepository node.Repository,
	redisInstance *configs.RedisInstance,
	protocolConfigurationRepository protocolConfiguration.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		pipelineRepository:              pipelineRepository,
		validatorService:                validatorService,
		dbConnection:                    dbConnection,
		viperConfig:                     viperConfig,
		pipelineNodeRepository:          pipelineNodeRepository,
		pipelineEdgeRepository:          pipelineEdgeRepository,
		nodeRepository:                  nodeRepository,
		redisInstance:                   redisInstance,
		protocolConfigurationRepository: protocolConfigurationRepository,
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

func (pipelineService *ServiceImpl) FindById(ginContext *gin.Context, pipelineId uint64) *model.PipelineResponse {
	var pipelineResponse *model.PipelineResponse
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntities, err := pipelineService.pipelineRepository.FindById(gormTransaction, pipelineId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		pipelineResponse = mapper.MapPipelineEntityIntoPipelineResponse(pipelineEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return pipelineResponse
}

func (pipelineService *ServiceImpl) RunPipeline(ginContext *gin.Context, pipelineId uint64) *model.PipelineResponse {
	var pipelineResponse *model.PipelineResponse
	err := pipelineService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		pipelineEntity, err := pipelineService.pipelineRepository.FindById(gormTransaction, pipelineId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		pipelineResponse = mapper.MapPipelineEntityIntoPipelineResponse(pipelineEntity)
		protocolConfigurationIdMap := map[uint64]struct{}{}
		for _, pipelineNodeResponse := range pipelineResponse.PipelineNode {
			id := pipelineNodeResponse.Config.ProtocolConfigurationId
			if id == 0 {
				continue
			}
			protocolConfigurationIdMap[id] = struct{}{}
		}

		// 2️⃣ Buat slice unik untuk query
		var protocolConfigurationIds []uint64
		for id := range protocolConfigurationIdMap {
			protocolConfigurationIds = append(protocolConfigurationIds, id)
		}
		protocolConfigurations, err := pipelineService.protocolConfigurationRepository.FindAllByIds(gormTransaction, protocolConfigurationIds)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		if len(protocolConfigurationIds) != len(protocolConfigurations) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, "Some protocol configuration not found"))
		}
		protocolConfigMap := make(map[uint64]model.ProtocolConfigurationResponse)
		for _, protocolConfigurationEntity := range protocolConfigurations {
			protocolConfigMap[protocolConfigurationEntity.Id] = *mapper.MapProtocolConfigurationEntityIntoProtocolConfigurationResponse(&protocolConfigurationEntity)
		}
		// 4️⃣ Mapping tiap node dengan ProtocolConfiguration-nya
		for _, pipelineNodeResponse := range pipelineResponse.PipelineNode {
			protocolConfigurationId := pipelineNodeResponse.Config.ProtocolConfigurationId
			if protocolConfigurationId == 0 {
				continue // tidak semua node harus punya protocol configuration
			}

			if protocolConfigurationResponse, ok := protocolConfigMap[protocolConfigurationId]; ok {
				pipelineNodeResponse.Config.ProtocolConfigurationResponse = protocolConfigurationResponse
			} else {
				// kalau tidak ditemukan, bisa log atau skip
				logrus.Debug("ProtocolConfiguration dengan ID %d tidak ditemukan\n", protocolConfigurationId)
			}
		}
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	// ✅ Publish ke Redis setelah data berhasil diambil
	if pipelineService.redisInstance != nil && pipelineResponse != nil {
		payload, err := json.Marshal(pipelineResponse)
		if err != nil {
			logrus.Debug("❌ Failed to marshal pipeline response: %v\n", err)
		} else {
			topic := fmt.Sprintf("pipeline")
			if err := pipelineService.redisInstance.Publish(topic, string(payload)); err != nil {
				logrus.Debug("❌ Failed to publish pipeline to Redis: %v\n", err)
			} else {
				logrus.Debug("📢 Published pipeline Id %d to topic '%s'\n", pipelineResponse.Id, topic)
			}
		}
	} else {
		logrus.Debug("⚠️ Redis instance or pipeline response is nil, skipping publish")
	}
	return pipelineResponse
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
		pipelineEntity, err := pipelineService.pipelineRepository.FindById(gormTransaction, updatePipelineRequest.Id)
		pipelineEntity.Auditable = entity.UpdateAuditable("System")
		err = pipelineService.pipelineEdgeRepository.DeleteByPipelineId(gormTransaction, pipelineEntity.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		err = pipelineService.pipelineNodeRepository.DeleteByPipelineId(gormTransaction, pipelineEntity.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		updatedPipelineEntity := mapper.MapUpdatePipelineRequestIntoPipelineEntity(updatePipelineRequest)
		updatedPipelineEntity.Auditable = entity.UpdateAuditable("System")
		err = pipelineService.pipelineRepository.Update(gormTransaction, updatedPipelineEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		mapOfTemporaryIdWithDatabaseId := map[string]uint64{}
		var nodeIds []uint64

		for i := range updatedPipelineEntity.PipelineNode {
			nodeIds = append(nodeIds, pipelineEntity.PipelineNode[i].NodeId)
		}

		nodeEntities, err := pipelineService.nodeRepository.FindBatchById(gormTransaction, nodeIds)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		if len(nodeEntities) != len(nodeIds) {
			return fmt.Errorf("not all nodes found: expected %d, got %d", len(nodeIds), len(nodeEntities))
		}

		// ✅ Fixed: Use index to work with pointers
		for i := range updatedPipelineEntity.PipelineNode {
			pipelineNodeEntity := pipelineEntity.PipelineNode[i]
			pipelineNodeEntity.PipelineId = pipelineEntity.Id

			err := pipelineService.pipelineNodeRepository.Create(gormTransaction, pipelineNodeEntity)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))

			mapOfTemporaryIdWithDatabaseId[pipelineNodeEntity.TempId] = pipelineNodeEntity.Id
		}

		// ✅ Fixed: Same for edges
		for i := range updatedPipelineEntity.PipelineEdge {
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
