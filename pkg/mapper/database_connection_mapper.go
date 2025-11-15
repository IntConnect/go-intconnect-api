package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapDatabaseConnectionEntityIntoDatabaseConnectionResponse(databaseConnectionEntity *entity.DatabaseConnection) *model.DatabaseConnectionResponse {
	var databaseConnectionResponse model.DatabaseConnectionResponse
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: helper.StringIntoTypeHookFunc,
		Result:     &databaseConnectionResponse,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))

	err = decoder.Decode(databaseConnectionEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	databaseConnectionResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&databaseConnectionEntity.Auditable)
	databaseConnectionResponse.Config = MapDatabaseConnectionConfigIntoDatabaseConnectionConfigResponse(databaseConnectionEntity.Config)
	return &databaseConnectionResponse
}

func MapDatabaseConnectionConfigIntoDatabaseConnectionConfigResponse(databaseConnectionConfig map[string]interface{}) *model.DatabaseConnectionConfigResponse {
	var databaseConnectionEntity model.DatabaseConnectionConfigResponse
	err := mapstructure.Decode(databaseConnectionConfig, &databaseConnectionEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &databaseConnectionEntity
}

func MapDatabaseConnectionEntitiesIntoDatabaseConnectionResponses(databaseConnectionEntities []entity.DatabaseConnection) []*model.DatabaseConnectionResponse {
	var databaseConnectionResponses []*model.DatabaseConnectionResponse
	for _, databaseConnectionEntity := range databaseConnectionEntities {
		databaseConnectionResponses = append(databaseConnectionResponses, MapDatabaseConnectionEntityIntoDatabaseConnectionResponse(&databaseConnectionEntity))
	}
	return databaseConnectionResponses
}

func MapCreateDatabaseConnectionRequestIntoDatabaseConnectionEntity(createDatabaseConnectionRequest *model.CreateDatabaseConnectionRequest) *entity.DatabaseConnection {
	var databaseConnectionEntity entity.DatabaseConnection
	err := mapstructure.Decode(createDatabaseConnectionRequest, &databaseConnectionEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &databaseConnectionEntity
}

func MapUpdateDatabaseConnectionRequestIntoDatabaseConnectionEntity(updateDatabaseConnectionRequest *model.UpdateDatabaseConnectionRequest, databaseConnectionEntity *entity.DatabaseConnection) {
	//helper.DecoderConfigMapper(updateDatabaseConnectionRequest, &databaseConnectionEntity)
}
