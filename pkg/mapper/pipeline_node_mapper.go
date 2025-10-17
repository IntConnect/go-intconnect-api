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

	// Paksa konversi jadi value, bukan pointer
	pipelineNodeEntity.NodeId = uint64(createPipelineNodeRequest.NodeID)
	pipelineNodeEntity.Type = createPipelineNodeRequest.Type
	pipelineNodeEntity.Label = createPipelineNodeRequest.Label
	pipelineNodeEntity.PositionX = createPipelineNodeRequest.PositionX
	pipelineNodeEntity.PositionY = createPipelineNodeRequest.PositionY
	pipelineNodeEntity.TempId = createPipelineNodeRequest.TempID
	pipelineNodeEntity.Config = createPipelineNodeRequest.Config

	// Optional jika ada config
	if createPipelineNodeRequest.Config != nil {
		jsonBytes, err := json.Marshal(createPipelineNodeRequest.Config)
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, "invalid config json", err))
		pipelineNodeEntity.ConfigRaw = jsonBytes
	}

	return &pipelineNodeEntity
}
