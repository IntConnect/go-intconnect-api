package machine_document

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (machineDocumentRepository *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.MachineDocument, error) {
	var machineEntities []entity.MachineDocument
	err := gormTransaction.Find(&machineEntities).Error
	return machineEntities, err
}

func (machineDocumentRepository *RepositoryImpl) FindById(gormTransaction *gorm.DB, machineId uint64) (*entity.MachineDocument, error) {
	var machineEntity entity.MachineDocument
	err := gormTransaction.Model(&entity.MachineDocument{}).Where("id = ?", machineId).Find(&machineEntity).Error
	return &machineEntity, err
}

func (machineDocumentRepository *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, machineIds []uint64) ([]*entity.MachineDocument, error) {
	var machineEntity []*entity.MachineDocument
	err := gormTransaction.Model(&entity.MachineDocument{}).Where("id IN ?", machineIds).Find(&machineEntity).Error
	return machineEntity, err
}

func (machineDocumentRepository *RepositoryImpl) Create(gormTransaction *gorm.DB, machineDocumentEntity *entity.MachineDocument) error {
	return gormTransaction.Create(machineDocumentEntity).Error
}

func (machineDocumentRepository *RepositoryImpl) CreateBatch(gormTransaction *gorm.DB, machineDocumentEntity []*entity.MachineDocument) error {
	return gormTransaction.Create(machineDocumentEntity).Error
}

func (machineDocumentRepository *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.MachineDocument{}).Where("id = ?", id).Delete(entity.MachineDocument{}).Error
}
