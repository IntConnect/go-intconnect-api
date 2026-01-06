package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/helper"
)

func FuncMapParameter(parameterEntity *entity.Parameter, parameterResponse *model.ParameterResponse) {
	processedParameterSequenceResponses := helper.MapEntitiesIntoResponses[*entity.ProcessedParameterSequence, model.ProcessedParameterSequenceResponse](parameterEntity.ProcessedParameterSequence)
	parameterResponse.ProcessedParameterSequence = processedParameterSequenceResponses
}
