package entity

type Role struct {
	Id          uint64        `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string        `gorm:"column:name"`
	Description string        `gorm:"column:description"`
	Permissions []*Permission `gorm:"many2many:roles_permissions;joinForeignKey:RoleId;joinReferences:PermissionId"`
	Auditable   Auditable     `gorm:"embedded"`
}

func (roleEntity *Role) GetAuditable() *Auditable {
	return &roleEntity.Auditable
}
