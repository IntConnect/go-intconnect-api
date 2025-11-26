package entity

type MachineDocument struct {
	Id          uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	MachineId   uint64 `gorm:"column:machine_id;"`
	Code        string `gorm:"column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	FilePath    string `gorm:"column:file_path"`
	Auditable
}
