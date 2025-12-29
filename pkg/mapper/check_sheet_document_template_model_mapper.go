package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"time"
)

func MapCheckSheetDocumentTemplate(checkSheetDocumentTemplateEntity *entity.CheckSheetDocumentTemplate, checkSheetDocumentTemplateResponse *model.CheckSheetDocumentTemplateResponse) {
	if checkSheetDocumentTemplateEntity == nil || checkSheetDocumentTemplateResponse == nil {
		return
	}
	checkSheetDocumentTemplateResponse.EffectiveDate = checkSheetDocumentTemplateEntity.EffectiveDate.Format(time.RFC3339)
}
