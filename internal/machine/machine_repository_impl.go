package machine

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (machineRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Machine, error) {
	var machineEntities []entity.Machine
	err := gormTransaction.Find(&machineEntities).Error
	return machineEntities, err
}

func (machineRepositoryImpl *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, machineIds []uint64) ([]entity.Machine, error) {
	var machineEntities []entity.Machine
	err := gormTransaction.Where("id IN ?", machineIds).Find(&machineEntities).Error
	return machineEntities, err
}

func (machineRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]entity.Machine, int64, error) {
	var machineEntities []entity.Machine
	var totalItems int64
	rawQuery := gormTransaction.Model(&entity.Machine{})
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}
	if err := rawQuery.
		Preload("Facility").
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&machineEntities).Error; err != nil {
		return nil, 0, err
	}

	return machineEntities, totalItems, nil
}

func (machineRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, machineId uint64) (*entity.Machine, error) {
	var machineEntity entity.Machine
	err := gormTransaction.Model(&entity.Machine{}).
		Preload("MachineDocuments").
		Preload("MqttTopic.Parameters").
		Preload("MqttTopic.MqttBroker").
		Where("id = ?", machineId).
		First(&machineEntity).Error

	return &machineEntity, err
}

func (machineRepositoryImpl *RepositoryImpl) FindByFacilityId(gormTransaction *gorm.DB, facilityId uint64) ([]*entity.Machine, error) {
	var machineEntities []*entity.Machine
	err := gormTransaction.Model(&entity.Machine{}).
		Preload("MachineDocuments").
		Preload("MqttTopic.Parameters").
		Preload("MqttTopic.MqttBroker").
		Where("facility_id = ?", facilityId).
		Find(&machineEntities).Error

	return machineEntities, err
}

func (machineRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Machine) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (machineRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, machineEntity *entity.Machine) error {
	return gormTransaction.Model(machineEntity).Updates(machineEntity).Error
}

func (machineRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, machineEntity *entity.Machine) error {
	return gormTransaction.Save(machineEntity).Error

}
