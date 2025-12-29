package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"time"
)

func MapCheckSheet(checkSheetEntity *entity.CheckSheet, checkSheetResponse *model.CheckSheetResponse) {
	if checkSheetEntity == nil || checkSheetResponse == nil {
		return
	}
	checkSheetResponse.Timestamp = checkSheetEntity.Timestamp.Format(time.RFC3339)
}
