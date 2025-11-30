package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

func MapFacilityEntityIntoFacilityResponse(facilityEntity *entity.Facility) *model.FacilityResponse {
	var facilityResponse *model.FacilityResponse
	facilityResponse = helper.DecodeFromSource[*entity.Facility, *model.FacilityResponse](facilityEntity, facilityResponse)
	facilityResponse.AuditableResponse = AuditableEntityIntoEntityResponse(&facilityEntity.Auditable)
	return facilityResponse
}

func MapFacilityEntitiesIntoFacilityResponses(facilityEntities []entity.Facility) []*model.FacilityResponse {
	var facilityResponses []*model.FacilityResponse
	for _, facilityEntity := range facilityEntities {
		facilityResponses = append(facilityResponses, MapFacilityEntityIntoFacilityResponse(&facilityEntity))
	}
	return facilityResponses
}

func MapCreateFacilityRequestIntoFacilityEntity(createFacilityRequest *model.CreateFacilityRequest) *entity.Facility {
	var facilityEntity entity.Facility
	err := mapstructure.Decode(createFacilityRequest, &facilityEntity)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
	return &facilityEntity
}

func MapUpdateFacilityRequestIntoFacilityEntity(updateFacilityRequest *model.UpdateFacilityRequest, facilityEntity *entity.Facility) {
	//helper.DecoderConfigMapper(updateFacilityRequest, &facilityEntity)
}
