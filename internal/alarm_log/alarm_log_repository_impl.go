package alarm_log

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (alarmLogRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.AlarmLog, error) {
	var alarmLogEntities []*entity.AlarmLog
	err := gormTransaction.Find(&alarmLogEntities).Error
	return alarmLogEntities, err
}

func (alarmLogRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.AlarmLog, int64, error) {

	var alarmLogEntities []*entity.AlarmLog
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.AlarmLog{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("action ILIKE ? OR feature ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Preload("AcknowledgedByUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Parameter", func(gormTransaction *gorm.DB) *gorm.DB {
			return gormTransaction.Select("id, name, code, machine_id").
				Preload("Machine", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, name")
				})
		}).
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&alarmLogEntities).Error; err != nil {
		return nil, 0, err
	}

	return alarmLogEntities, totalItems, nil
}

func (alarmLogRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, alarmLogId uint64) (*entity.AlarmLog, error) {
	var alarmLogEntity entity.AlarmLog
	err := gormTransaction.Model(&entity.AlarmLog{}).Where("id = ?", alarmLogId).Find(&alarmLogEntity).Error

	return &alarmLogEntity, err
}

func (alarmLogRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, alarmLogEntity *entity.AlarmLog) error {
	return gormTransaction.Model(alarmLogEntity).Create(alarmLogEntity).Error

}

func (alarmLogRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, alarmLogEntity *entity.AlarmLog) error {
	return gormTransaction.Model(alarmLogEntity).Updates(alarmLogEntity).Error

}
