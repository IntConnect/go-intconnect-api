package entity

type Permission struct {
	Id          uint64    `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	Category    string    `gorm:"column:category"`
	Description string    `gorm:"column:description"`
	Roles       []*Role   `gorm:"many2many:roles_permissions;joinForeignKey:PermissionId;joinReferences:RoleId"`
	Auditable   Auditable `gorm:"embedded"`
}

func (permissionEntity *Permission) GetAuditable() *Auditable {
	return &permissionEntity.Auditable
}
