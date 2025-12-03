package entity

type RolePermission struct {
	Id           uint64 `gorm:"column:id"`
	RoleId       uint64 `gorm:"column:role_id"`
	PermissionId uint64 `gorm:"column:permission_id"`
	Granted      bool   `gorm:"column:granted;default:true"`
}

func (RolePermission) TableName() string {
	return "roles_permissions"
}
