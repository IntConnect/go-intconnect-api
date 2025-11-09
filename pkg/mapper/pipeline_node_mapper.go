package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapCreatePipelineNodeRequestsIntoPipelineNodeEntities(createPipelineNodeRequests []model.CreatePipelineNodeRequest) []*entity.PipelineNode {
	var pipelineNodeEntities []*entity.PipelineNode
	for _, createPipelineNodeRequest := range createPipelineNodeRequests {
		pipelineNodeEntities = append(pipelineNodeEntities, MapCreatePipelineNodeRequestIntoPipelineNodeEntity(createPipelineNodeRequest))
	}
	return pipelineNodeEntities
}

func MapCreatePipelineNodeRequestIntoPipelineNodeEntity(createPipelineNodeRequest model.CreatePipelineNodeRequest) *entity.PipelineNode {
	var pipelineNodeEntity entity.PipelineNode
	err := mapstructure.Decode(createPipelineNodeRequest, &pipelineNodeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineNodeEntity.Config = createPipelineNodeRequest.Config
	return &pipelineNodeEntity
}

func MapPipelineNodeEntitiesIntoPipelineNodeResponse(pipelineNodeEntities []*entity.PipelineNode) []*model.PipelineNodeResponse {
	var pipelineNodeResponses []*model.PipelineNodeResponse
	for _, pipelineNodeEntity := range pipelineNodeEntities {
		pipelineNodeResponses = append(pipelineNodeResponses, MapPipelineNodeEntityIntoPipelineNodeResponse(pipelineNodeEntity))
	}
	return pipelineNodeResponses
}

func MapPipelineNodeEntityIntoPipelineNodeResponse(pipelineNodeEntity *entity.PipelineNode) *model.PipelineNodeResponse {
	var pipelineNodeResponse model.PipelineNodeResponse
	err := mapstructure.Decode(pipelineNodeEntity, &pipelineNodeResponse)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineNodeResponse.Config = MapPipelineNodeConfigEntityIntoPipelineNodeConfigResponse(pipelineNodeEntity.Config)
	pipelineNodeResponse.Appearance = MapPipelineNodeAppearanceEntityIntoPipelineNodeAppearanceResponse(pipelineNodeEntity.Appearance)
	pipelineNodeResponse.NodeResponse = MapNodeEntityIntoNodeResponse(&pipelineNodeEntity.Node)
	return &pipelineNodeResponse
}

func MapPipelineNodeConfigEntityIntoPipelineNodeConfigResponse(pipelineNodeConfig map[string]interface{}) *model.PipelineNodeConfig {
	var parsedPipelineNodeConfig *model.PipelineNodeConfig
	err := mapstructure.Decode(pipelineNodeConfig, &parsedPipelineNodeConfig)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return parsedPipelineNodeConfig
}

func MapPipelineNodeAppearanceEntityIntoPipelineNodeAppearanceResponse(pipelineNodeAppearance map[string]interface{}) *model.PipelineNodeAppearance {
	var parsedPipelineNodeAppearance *model.PipelineNodeAppearance
	err := mapstructure.Decode(pipelineNodeAppearance, &parsedPipelineNodeAppearance)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return parsedPipelineNodeAppearance
}
