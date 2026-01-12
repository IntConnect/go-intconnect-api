package entity

type RolePermission struct {
	RoleId       uint64 `gorm:"column:role_id"`
	PermissionId uint64 `gorm:"column:permission_id"`
}

func (RolePermission) TableName() string {
	return "roles_permissions"
}
