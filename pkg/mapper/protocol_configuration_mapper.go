package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapProtocolConfigurationEntityIntoProtocolConfigurationResponse(protocolConfigurationEntity *entity.ProtocolConfiguration) *model.ProtocolConfigurationResponse {
	var protocolConfigurationResponse model.ProtocolConfigurationResponse
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: helper.StringIntoTypeHookFunc,
		Result:     &protocolConfigurationResponse,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	protocolConfigurationResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&protocolConfigurationEntity.Auditable)
	err = decoder.Decode(protocolConfigurationEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &protocolConfigurationResponse
}

func MapProtocolConfigurationEntitiesIntoProtocolConfigurationResponses(protocolConfigurationEntities []entity.ProtocolConfiguration) []*model.ProtocolConfigurationResponse {
	var protocolConfigurationResponses []*model.ProtocolConfigurationResponse
	for _, protocolConfigurationEntity := range protocolConfigurationEntities {
		protocolConfigurationResponses = append(protocolConfigurationResponses, MapProtocolConfigurationEntityIntoProtocolConfigurationResponse(&protocolConfigurationEntity))
	}
	return protocolConfigurationResponses
}

func MapCreateProtocolConfigurationRequestIntoProtocolConfigurationEntity(createProtocolConfigurationRequest *model.CreateProtocolConfigurationRequest) *entity.ProtocolConfiguration {
	var protocolConfigurationEntity entity.ProtocolConfiguration
	err := mapstructure.Decode(createProtocolConfigurationRequest, &protocolConfigurationEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &protocolConfigurationEntity
}

func MapUpdateProtocolConfigurationRequestIntoProtocolConfigurationEntity(updateProtocolConfigurationRequest *model.UpdateProtocolConfigurationRequest, protocolConfigurationEntity *entity.ProtocolConfiguration) {
	//helper.DecoderConfigMapper(updateProtocolConfigurationRequest, &protocolConfigurationEntity)
}
