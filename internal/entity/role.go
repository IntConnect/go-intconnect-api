package entity

type Role struct {
	Id           uint64 `gorm:"column:id,primaryKey,autoIncrement"`
	Name         string `gorm:"column:name"`
	Description  string `gorm:"column:description"`
	IsSystemRole bool   `gorm:"column:is_system_role"`
	Auditable
}
