package mapper

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapNodeEntityIntoNodeResponse(nodeEntity *entity.Node) *model.NodeResponse {
	var nodeResponse *model.NodeResponse
	nodeResponse = helper.DecodeFromSource[*model.NodeResponse](nodeEntity, nodeResponse)
	nodeResponse.DefaultConfig = nodeEntity.DefaultConfig
	nodeResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&nodeEntity.Auditable)
	return nodeResponse
}

func MapNodeEntitiesIntoNodeResponses(nodeEntities []entity.Node) []*model.NodeResponse {
	var nodeResponses []*model.NodeResponse
	for _, nodeEntity := range nodeEntities {
		fmt.Println(nodeEntity.DefaultConfig)
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
