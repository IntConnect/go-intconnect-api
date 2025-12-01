package seeders

import (
	"fmt"
	"go-intconnect-api/internal/entity"
	"strings"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
)

type PermissionSeeder struct{}

func (permissionSeeder *PermissionSeeder) Run(gormDatabase *gorm.DB) error {
	fakerInstance := faker.New()

	categories := []string{"machine", "facility", "user", "auth", "role"}

	for i := 0; i < 10; i++ {
		randomWord := fakerInstance.Lorem().Word()
		randomCategory := fakerInstance.RandomStringElement(categories)

		permissionEntity := entity.Permission{
			Code:        fmt.Sprintf("perm_%s_%d", randomCategory, fakerInstance.Int64Between(1, 9999)),
			Name:        strings.Title(fmt.Sprintf("%s %s", randomCategory, randomWord)),
			Category:    randomCategory,
			Description: fakerInstance.Lorem().Sentence(10),
			Auditable:   entity.NewAuditable("Admin"),
		}

		if err := gormDatabase.Create(&permissionEntity).Error; err != nil {
			return err
		}

	}

	return nil
}
