package seeders

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go-intconnect-api/internal/entity"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
)

type MachineSeeder struct{}

func (machineSeeder *MachineSeeder) Run(gormDatabase *gorm.DB) error {
	// reset table
	gormDatabase.Exec("TRUNCATE TABLE machines RESTART IDENTITY CASCADE")

	fake := faker.New()
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 15; i++ {
		machineName := fmt.Sprintf("Machine %d - %s", i, fake.Person().Name())
		machine := entity.Machine{
			FacilityId:       1,
			Name:             machineName,
			Code:             strings.ToLower(machineName),
			Description:      "",
			CameraX:          fake.Float(2, -50, 50),
			CameraY:          fake.Float(2, 5, 30),
			CameraZ:          fake.Float(2, -50, 50),
			ThumbnailPath:    fmt.Sprintf("machines/photos/machine_%d.png", i),
			ModelPath:        fmt.Sprintf("machines/models/machine_%d.glb", i),
			MachineDocuments: nil,
		}

		if err := gormDatabase.Model(&entity.Machine{}).Create(&machine).Error; err != nil {
			return err
		}
	}

	return nil
}
