package mapper

import (
	"encoding/json"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapCreatePipelineNodeDtosIntoPipelineNodeEntities(createPipelineNodeDtos []model.CreatePipelineNodeDto) []*entity.PipelineNode {
	var pipelineNodeEntities []*entity.PipelineNode
	for _, createPipelineNodeDto := range createPipelineNodeDtos {
		pipelineNodeEntities = append(pipelineNodeEntities, MapCreatePipelineNodeDtoIntoPipelineNodeEntity(createPipelineNodeDto))
	}
	return pipelineNodeEntities
}

func MapCreatePipelineNodeDtoIntoPipelineNodeEntity(createPipelineNodeDto model.CreatePipelineNodeDto) *entity.PipelineNode {
	var pipelineNodeEntity entity.PipelineNode
	err := mapstructure.Decode(createPipelineNodeDto, &pipelineNodeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	// Paksa konversi jadi value, bukan pointer
	pipelineNodeEntity.NodeId = uint64(createPipelineNodeDto.NodeID)
	pipelineNodeEntity.Type = createPipelineNodeDto.Type
	pipelineNodeEntity.Label = createPipelineNodeDto.Label
	pipelineNodeEntity.PositionX = createPipelineNodeDto.PositionX
	pipelineNodeEntity.PositionY = createPipelineNodeDto.PositionY
	pipelineNodeEntity.TempId = createPipelineNodeDto.TempID
	pipelineNodeEntity.Config = createPipelineNodeDto.Config

	// Optional jika ada config
	if createPipelineNodeDto.Config != nil {
		jsonBytes, err := json.Marshal(createPipelineNodeDto.Config)
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, "invalid config json", err))
		pipelineNodeEntity.ConfigRaw = jsonBytes
	}

	return &pipelineNodeEntity
}
