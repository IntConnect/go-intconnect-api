package entity

import "go-intconnect-api/internal/trait"

type User struct {
	Id         uint64           `gorm:"column:id;primaryKey;autoIncrement"`
	RoleId     uint64           `gorm:"column:role_id;"`
	Username   string           `gorm:"column:username"`
	Name       string           `gorm:"column:name"`
	Email      string           `gorm:"column:email"`
	Password   string           `gorm:"column:password"`
	AvatarPath string           `gorm:"column:avatar_path"`
	Status     trait.UserStatus `gorm:"column:status"`
	Role       *Role            `gorm:"foreignKey:RoleId;references:Id"`
	AuditLog   []*AuditLog      `gorm:"foreignKey:UserId;references:Id"`
	Auditable
}

func (userEntity User) GetAuditable() *Auditable {
	return &userEntity.Auditable
}
