package seeders

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RolePermissionSeeder struct{}

func (rolePermissionSeeder *RolePermissionSeeder) Run(gormDatabase *gorm.DB) error {
	var permissionEntities []entity.Permission
	err := gormDatabase.Model(&entity.Permission{}).Find(&permissionEntities).Error
	if err != nil {
		panic(err)
	}
	for _, permissionEntity := range permissionEntities {
		err = gormDatabase.Create(&entity.RolePermission{
			RoleId:       1,
			PermissionId: permissionEntity.Id,
			Granted:      false,
		}).Error
		if err != nil {
			panic(err)
		}
	}
	return nil
}
