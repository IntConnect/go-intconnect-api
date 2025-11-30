package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
)

func FuncMapAuditable[S entity.HasAuditable, R model.HasAuditableResponse](
	entityObject S,
	responseObject R,
) {
	responseObject.SetAuditableResponse(
		AuditableEntityIntoEntityResponse(entityObject.GetAuditable()),
	)
}

func FuncMapSimpleAuditable[S entity.HasSimpleAuditable, R model.HasSimpleAuditableResponse](
	entityObject S,
	responseObject R,
) {
	responseObject.SetSimpleAuditableResponse(
		SimpleAuditableEntityIntoSimpleEntityResponse(entityObject.GetSimpleAuditable()),
	)
}
