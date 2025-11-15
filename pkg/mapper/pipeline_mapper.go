package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapPipelineEntityIntoPipelineResponse(pipelineEntity *entity.Pipeline) *model.PipelineResponse {
	var pipelineResponse *model.PipelineResponse
	pipelineResponse = helper.DecodeFromSource[*entity.Pipeline, *model.PipelineResponse](pipelineEntity, pipelineResponse)
	pipelineResponse.PipelineNode = MapPipelineNodeEntitiesIntoPipelineNodeResponse(pipelineEntity.PipelineNode)
	pipelineResponse.PipelineEdge = MapPipelineEdgeEntitiesIntoPipelineEdgeResponse(pipelineEntity.PipelineEdge)
	pipelineResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&pipelineEntity.Auditable)

	return pipelineResponse
}

func MapPipelineEntitiesIntoPipelineResponses(pipelineEntities []*entity.Pipeline) []*model.PipelineResponse {
	var pipelineResponses []*model.PipelineResponse
	for _, pipelineEntity := range pipelineEntities {
		pipelineResponses = append(pipelineResponses, MapPipelineEntityIntoPipelineResponse(pipelineEntity))
	}
	return pipelineResponses
}

func MapCreatePipelineRequestIntoPipelineEntity(createPipelineRequest *model.CreatePipelineRequest) *entity.Pipeline {
	var pipelineEntity entity.Pipeline
	err := mapstructure.Decode(createPipelineRequest, &pipelineEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineEntity.PipelineNode = MapCreatePipelineNodeRequestsIntoPipelineNodeEntities(createPipelineRequest.Nodes)
	pipelineEntity.PipelineEdge = MapCreatePipelineEdgeRequestsIntoPipelineEdgeEntities(createPipelineRequest.Edges)
	pipelineEntity.Config = createPipelineRequest.Config
	return &pipelineEntity
}

func MapUpdatePipelineRequestIntoPipelineEntity(createPipelineRequest *model.UpdatePipelineRequest) *entity.Pipeline {
	var pipelineEntity entity.Pipeline
	err := mapstructure.Decode(createPipelineRequest, &pipelineEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineEntity.PipelineNode = MapCreatePipelineNodeRequestsIntoPipelineNodeEntities(createPipelineRequest.Nodes)
	pipelineEntity.PipelineEdge = MapCreatePipelineEdgeRequestsIntoPipelineEdgeEntities(createPipelineRequest.Edges)
	pipelineEntity.Config = createPipelineRequest.Config
	return &pipelineEntity
}
