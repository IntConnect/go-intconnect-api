package parameter

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

func (machineRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Machine, int64, error) {
	var machineEntities []entity.Machine
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("machinename LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.Machine{}).
		Preload("MachineGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&machineEntities).Error
	gormTransaction.Model(&entity.Machine{}).Count(&totalItems)
	return machineEntities, totalItems, err
}

func (machineRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, machineId uint64) (*entity.Machine, error) {
	var machineEntity entity.Machine
	err := gormTransaction.Model(&entity.Machine{}).
		Preload("MachineGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", machineId).Find(&machineEntity).Error

	return &machineEntity, err
}

func (machineRepositoryImpl *RepositoryImpl) FindByName(pipelineName string) *entity.Machine {
	//TODO implement me
	panic("implement me")
}

func (machineRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Machine) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (machineRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, machineEntity *entity.Machine) error {
	return gormTransaction.Model(machineEntity).Save(machineEntity).Error
}

func (machineRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Machine{}).Where("id = ?", id).Delete(entity.Machine{}).Error
}
