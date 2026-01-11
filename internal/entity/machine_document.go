package entity

type MachineDocument struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId   uint64    `gorm:"column:machine_id;"`
	Name        string    `gorm:"column:name"`
	FilePath    string    `gorm:"column:file_path"`
	Description string    `gorm:"column:description"`
	Machine     *Machine  `gorm:"foreignKey:MachineId;references:Id"`
	Auditable   Auditable `gorm:"embedded"`
}

func (machineDocumentEntity *MachineDocument) GetId() uint64 {
	return machineDocumentEntity.Id
}

func (machineDocumentEntity *MachineDocument) GetAuditable() *Auditable {
	return &machineDocumentEntity.Auditable
}
