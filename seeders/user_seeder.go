package seeders

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/trait"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (userSeeder *UserSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	gormDatabase.Model(&entity.User{}).Create(&entity.User{
		RoleId:     1,
		Username:   "admin",
		Name:       "Administrator",
		Email:      "admin@gmail.com",
		Password:   string(hashedPassword),
		AvatarPath: "",
		Status:     trait.UserStatusActive,
		Auditable:  entity.NewAuditable("Administrator"),
	})
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("mngr"), bcrypt.DefaultCost)

	gormDatabase.Model(&entity.User{}).Create(&entity.User{
		RoleId:     2,
		Username:   "mngr",
		Name:       "Manager",
		Email:      "mngr@gmail.com",
		Password:   string(hashedPassword),
		AvatarPath: "",
		Status:     trait.UserStatusActive,
		Auditable:  entity.NewAuditable("Administrator"),
	})
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("spv"), bcrypt.DefaultCost)

	gormDatabase.Model(&entity.User{}).Create(&entity.User{
		RoleId:     3,
		Username:   "spv",
		Name:       "Supervisor",
		Email:      "supervisor@gmail.com",
		Password:   string(hashedPassword),
		AvatarPath: "",
		Status:     trait.UserStatusActive,
		Auditable:  entity.NewAuditable("Administrator"),
	})
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("mtn"), bcrypt.DefaultCost)

	gormDatabase.Model(&entity.User{}).Create(&entity.User{
		RoleId:     4,
		Username:   "mtn",
		Name:       "Maintenance",
		Email:      "maintenance@gmail.com",
		Password:   string(hashedPassword),
		AvatarPath: "",
		Status:     trait.UserStatusActive,
		Auditable:  entity.NewAuditable("Administrator"),
	})
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("oprt"), bcrypt.DefaultCost)

	gormDatabase.Model(&entity.User{}).Create(&entity.User{
		RoleId:     5,
		Username:   "oprt",
		Name:       "Operator",
		Email:      "operator@gmail.com",
		Password:   string(hashedPassword),
		AvatarPath: "",
		Status:     trait.UserStatusActive,
		Auditable:  entity.NewAuditable("Administrator"),
	})
	return nil
}
