package machine_document

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.MachineDocument, error)
	FindById(gormTransaction *gorm.DB, machineId uint64) (*entity.MachineDocument, error)
	FindBatchById(gormTransaction *gorm.DB, machineIds []uint64) ([]*entity.MachineDocument, error)
	Create(gormTransaction *gorm.DB, machineDocumentEntity *entity.MachineDocument) error
	CreateBatch(gormTransaction *gorm.DB, machineDocumentEntities []*entity.MachineDocument) error
	UpdateBatch(gormTransaction *gorm.DB, machineDocumentEntities []*entity.MachineDocument) error
	Delete(gormTransaction *gorm.DB, machineDocumentId uint64) error
}
