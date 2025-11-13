package entity

type Permission struct {
	Id          uint64 `gorm:"primaryKey;auto_increment"`
	Code        string `gorm:"column:code"`
	Name        string `gorm:"column:name"`
	Category    string `gorm:"column:category"`
	Description string `gorm:"column:description"`
	Auditable
}
