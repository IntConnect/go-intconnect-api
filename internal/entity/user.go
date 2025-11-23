package entity

import "go-intconnect-api/internal/trait"

type User struct {
	Id       uint64           `gorm:"column:id;primaryKey;autoIncrement"`
	Username string           `gorm:"column:username"`
	Name     string           `gorm:"column:name"`
	Email    string           `gorm:"column:email"`
	Password string           `gorm:"column:password"`
	Avatar   string           `gorm:"column:avatar"`
	Status   trait.UserStatus `gorm:"column:status"`
	Auditable
}

func (userEntity User) GetAuditable() *Auditable {
	return &userEntity.Auditable
}
