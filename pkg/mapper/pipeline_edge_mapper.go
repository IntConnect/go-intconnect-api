package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapCreatePipelineEdgeDtosIntoPipelineEdgeEntities(createPipelineEdgeDtos []model.CreatePipelineEdgeDto) []*entity.PipelineEdge {
	var pipelineEdgeEntities []*entity.PipelineEdge
	for _, createPipelineEdgeDto := range createPipelineEdgeDtos {
		pipelineEdgeEntities = append(pipelineEdgeEntities, MapCreatePipelineEdgeDtoIntoPipelineEdgeEntity(createPipelineEdgeDto))
	}
	return pipelineEdgeEntities
}

func MapCreatePipelineEdgeDtoIntoPipelineEdgeEntity(createPipelineEdgeDto model.CreatePipelineEdgeDto) *entity.PipelineEdge {
	var pipelineEdgeEntity entity.PipelineEdge
	err := mapstructure.Decode(createPipelineEdgeDto, &pipelineEdgeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &pipelineEdgeEntity
}
