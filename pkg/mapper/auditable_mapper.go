package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/helper"
)

func AuditableEntityIntoEntityResponse(auditableEntity *entity.Auditable) *model.AuditableResponse {
	var auditableResponse model.AuditableResponse
	auditableResponse.CreatedBy = auditableEntity.CreatedBy
	auditableResponse.CreatedAt = auditableEntity.CreatedAt.String()
	auditableResponse.UpdatedBy = auditableEntity.UpdatedBy
	auditableResponse.UpdatedAt = auditableEntity.UpdatedAt.String()
	helper.CheckPointerWrapper(auditableEntity.DeletedBy, func() {
		auditableResponse.DeletedBy = *auditableEntity.DeletedBy
	})
	auditableResponse.DeletedBy = auditableEntity.DeletedAt.Time.String()
	return &auditableResponse

}
