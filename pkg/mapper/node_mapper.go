package mapper

import (
	"github.com/go-viper/mapstructure/v2"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
)

func MapNodeEntityIntoNodeResponse(nodeEntity *entity.Node) *model.NodeResponse {
	var nodeResponse model.NodeResponse
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: helper.StringIntoTypeHookFunc,
		Result:     &nodeResponse,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	nodeResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&nodeEntity.Auditable)
	err = decoder.Decode(nodeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &nodeResponse
}

func MapNodeEntitiesIntoNodeResponses(nodeEntities []entity.Node) []*model.NodeResponse {
	var nodeResponses []*model.NodeResponse
	for _, nodeEntity := range nodeEntities {
		nodeResponses = append(nodeResponses, MapNodeEntityIntoNodeResponse(&nodeEntity))
	}
	return nodeResponses
}

func MapCreateNodeRequestIntoNodeEntity(createNodeRequest *model.CreateNodeRequest) *entity.Node {
	var nodeEntity entity.Node
	err := mapstructure.Decode(createNodeRequest, &nodeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &nodeEntity
}

func MapUpdateNodeRequestIntoNodeEntity(updateNodeRequest *model.UpdateNodeRequest, nodeEntity *entity.Node) {
	helper.DecoderConfigMapper(updateNodeRequest, &nodeEntity)
}
