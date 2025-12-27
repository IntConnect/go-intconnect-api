package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/trait"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
)

type FacilitySeeder struct{}

func (facilitySeeder *FacilitySeeder) Run(gormDatabase *gorm.DB) error {
	// reset table
	gormDatabase.Exec("TRUNCATE TABLE facilities RESTART IDENTITY CASCADE")

	fake := faker.New()
	rand.Seed(time.Now().UnixNano())

	statuses := []trait.FacilityStatus{
		trait.FacilityStatusActive,
		trait.FacilityStatusMaintenance,
		trait.FacilityStatusArchived,
	}

	for i := 1; i <= 15; i++ {
		facility := entity.Facility{
			Name:        fmt.Sprintf("Facility %s", fake.Company().Name()),
			Code:        fmt.Sprintf("FAC-%03d", i),
			Location:    fake.Address().City(),
			Description: fake.Lorem().Sentence(10),
			Status:      statuses[rand.Intn(len(statuses))],

			ThumbnailPath: fmt.Sprintf("thumbnails/facility_%d.png", i),
			ModelPath:     fmt.Sprintf("models/facility_%d.glb", i),

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
	}

	return nil
}
