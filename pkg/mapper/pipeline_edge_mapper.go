package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapCreatePipelineEdgeRequestsIntoPipelineEdgeEntities(createPipelineEdgeRequests []model.CreatePipelineEdgeRequest) []*entity.PipelineEdge {
	var pipelineEdgeEntities []*entity.PipelineEdge
	for _, createPipelineEdgeRequest := range createPipelineEdgeRequests {
		pipelineEdgeEntities = append(pipelineEdgeEntities, MapCreatePipelineEdgeRequestIntoPipelineEdgeEntity(createPipelineEdgeRequest))
	}
	return pipelineEdgeEntities
}

func MapCreatePipelineEdgeRequestIntoPipelineEdgeEntity(createPipelineEdgeRequest model.CreatePipelineEdgeRequest) *entity.PipelineEdge {
	var pipelineEdgeEntity entity.PipelineEdge
	err := mapstructure.Decode(createPipelineEdgeRequest, &pipelineEdgeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &pipelineEdgeEntity
}
