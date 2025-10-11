package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapProtocolConfigurationEntityIntoProtocolConfigurationResponse(nodeEntity *entity.ProtocolConfiguration) *model.ProtocolConfigurationResponse {
	var nodeResponse model.ProtocolConfigurationResponse
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

func MapProtocolConfigurationEntitiesIntoProtocolConfigurationResponses(nodeEntities []entity.ProtocolConfiguration) []*model.ProtocolConfigurationResponse {
	var nodeResponses []*model.ProtocolConfigurationResponse
	for _, nodeEntity := range nodeEntities {
		nodeResponses = append(nodeResponses, MapProtocolConfigurationEntityIntoProtocolConfigurationResponse(&nodeEntity))
	}
	return nodeResponses
}

func MapCreateProtocolConfigurationDtoIntoProtocolConfigurationEntity(createProtocolConfigurationDto *model.CreateProtocolConfigurationDto) *entity.ProtocolConfiguration {
	var nodeEntity entity.ProtocolConfiguration
	err := mapstructure.Decode(createProtocolConfigurationDto, &nodeEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &nodeEntity
}

func MapUpdateProtocolConfigurationDtoIntoProtocolConfigurationEntity(updateProtocolConfigurationDto *model.UpdateProtocolConfigurationDto, nodeEntity *entity.ProtocolConfiguration) {
	helper.DecoderConfigMapper(updateProtocolConfigurationDto, &nodeEntity)
}
