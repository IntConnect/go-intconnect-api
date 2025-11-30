package seeders

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RoleSeeder struct{}

func (roleSeeder *RoleSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Create(&entity.Role{
		Name:        "Administrator",
		Description: "administrator",
		Auditable:   entity.NewAuditable("Administrator"),
	})
	gormDatabase.Create(&entity.Role{
		Name:        "Manager",
		Description: "manager",
		Auditable:   entity.NewAuditable("Manager"),
	})
	gormDatabase.Create(&entity.Role{
		Name:        "Supervisor",
		Description: "supervisor",
		Auditable:   entity.NewAuditable("Supervisor"),
	})
	gormDatabase.Create(&entity.Role{
		Name:        "Maintenance",
		Description: "maintenance",
		Auditable:   entity.NewAuditable("Maintenance"),
	})
	gormDatabase.Create(&entity.Role{
		Name:        "Operator",
		Description: "operator",
		Auditable:   entity.NewAuditable("Operator"),
	})

	return nil
}
