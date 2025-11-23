package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
)

func FuncMapAuditable[S entity.HasAuditable, R model.HasAuditableResponse](
	entityObject S,
	responseObject R, // ✅ R sudah pointer type
) {
	responseObject.SetAuditableResponse(
		AuditableEntityIntoEntityResponse(entityObject.GetAuditable()),
	)
}
