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
	pipelineResponse = helper.DecodeFromSource[*model.PipelineResponse](pipelineEntity, pipelineResponse)
	pipelineResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&pipelineEntity.Auditable)
	return pipelineResponse
}

func MapPipelineEntitiesIntoPipelineResponses(pipelineEntities []entity.Pipeline) []*model.PipelineResponse {
	var pipelineResponses []*model.PipelineResponse
	for _, pipelineEntity := range pipelineEntities {
		pipelineResponses = append(pipelineResponses, MapPipelineEntityIntoPipelineResponse(&pipelineEntity))
	}
	return pipelineResponses
}

func MapCreatePipelineDtoIntoPipelineEntity(createPipelineDto *model.CreatePipelineDto) *entity.Pipeline {
	var pipelineEntity entity.Pipeline
	err := mapstructure.Decode(createPipelineDto, &pipelineEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	pipelineEntity.PipelineNode = MapCreatePipelineNodeDtosIntoPipelineNodeEntities(createPipelineDto.Nodes)
	pipelineEntity.PipelineEdge = MapCreatePipelineEdgeDtosIntoPipelineEdgeEntities(createPipelineDto.Edges)
	pipelineEntity.Config = createPipelineDto.Config
	return &pipelineEntity
}

func MapUpdatePipelineDtoIntoPipelineEntity(updatePipelineDto *model.UpdatePipelineDto, pipelineEntity *entity.Pipeline) {
	helper.DecoderConfigMapper(updatePipelineDto, &pipelineEntity)
}
