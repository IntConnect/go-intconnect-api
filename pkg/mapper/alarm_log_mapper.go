package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"time"
)

func FuncMapAlarmLog(alarmLogEntity *entity.AlarmLog, alarmLogResponse *model.AlarmLogResponse) {
	alarmLogResponse.CreatedAt = alarmLogEntity.CreatedAt.Format(time.RFC3339)
	alarmLogResponse.UpdatedAt = alarmLogEntity.UpdatedAt.Format(time.RFC3339)
	alarmLogResponse.AcknowledgedAt = alarmLogEntity.AcknowledgedAt.Format(time.RFC3339)
}
