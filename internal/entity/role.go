package entity

type Role struct {
	Id          uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Auditable
}

func (roleEntity Role) GetAuditable() *Auditable {
	return &roleEntity.Auditable
}
