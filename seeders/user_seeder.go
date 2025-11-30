package seeders

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/trait"
	"strings"

	"github.com/jaswdr/faker/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (userSeeder *UserSeeder) Run(gormDatabase *gorm.DB) error {
	fakerInstance := faker.New()
	for i := 0; i < 10; i++ {
		// Generate fake data
		fullName := fakerInstance.Person().Name()
		firstName := fakerInstance.Person().FirstName()
		lastName := fakerInstance.Person().LastName()

		username := fmt.Sprintf("%s_%s", strings.ToLower(firstName), strings.ToLower(lastName))
		email := fakerInstance.Internet().Email()

		// Hash default password "password123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		userEntity := entity.User{
			Username:   username,
			Name:       fullName,
			Email:      email,
			Password:   string(hashedPassword),
			AvatarPath: "",
			RoleId:     1,
			Status:     trait.UserStatusActive,
			Auditable:  entity.NewAuditable("Seeder"),
		}

		// Insert into DB
		if err := gormDatabase.Create(&userEntity).Error; err != nil {
			return err
		}

		fmt.Printf("Inserted user: %s (%s)\n", userEntity.Name, userEntity.Email)
	}

	return nil
}
