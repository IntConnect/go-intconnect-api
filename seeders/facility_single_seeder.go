package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"go-intconnect-api/internal/entity"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
)

type FacilitySingleSeeder struct{}

func (facilitySeeder *FacilitySingleSeeder) Run(gormDatabase *gorm.DB) error {
	// reset table
	gormDatabase.Exec("TRUNCATE TABLE facilities RESTART IDENTITY CASCADE")

	fake := faker.New()
	rand.Seed(time.Now().UnixNano())

	facility := entity.Facility{
		Name:        fmt.Sprintf("Facility HVAC"),
		Code:        fmt.Sprintf("FAC-1"),
		Location:    fake.Address().City(),
		Description: fake.Lorem().Sentence(10),

		ThumbnailPath: fmt.Sprintf("thumbnails/lanjutkan.png"),
		ModelPath:     fmt.Sprintf("models/full-kmi.glb"),

		PositionX: fake.Float64(2, -50, 50),
		PositionY: fake.Float64(2, 0, 20),
		PositionZ: fake.Float(2, -50, 50),

		CameraX: fake.Float(2, -100, 100),
		CameraY: fake.Float(2, 10, 50),
		CameraZ: fake.Float(2, -100, 100),

		Auditable: entity.NewAuditable("Seeder"),
	}

	if err := gormDatabase.Model(&entity.Facility{}).Create(&facility).Error; err != nil {
		return err
	}

	return nil
}
