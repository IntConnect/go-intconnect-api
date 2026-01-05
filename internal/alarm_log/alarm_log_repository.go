package alarm_log

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.AlarmLog, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.AlarmLog, int64, error)
	FindById(gormTransaction *gorm.DB, alarmLogId uint64) (*entity.AlarmLog, error)
	Create(gormTransaction *gorm.DB, alarmLogEntity *entity.AlarmLog) error
	Update(gormTransaction *gorm.DB, alarmLogEntity *entity.AlarmLog) error
}
