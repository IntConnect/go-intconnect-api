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

type MachineSingleSeeder struct{}

func (machineSeeder *MachineSingleSeeder) Run(gormDatabase *gorm.DB) error {
	// reset table
	gormDatabase.Exec("TRUNCATE TABLE machines RESTART IDENTITY CASCADE")

	fake := faker.New()
	rand.Seed(time.Now().UnixNano())

	machineName := fmt.Sprintf("Chiller AICOOL TR300")
	machine := entity.Machine{
		FacilityId:       1,
		Name:             machineName,
		Code:             strings.ToLower(machineName),
		Description:      "",
		CameraX:          fake.Float(2, -50, 50),
		CameraY:          fake.Float(2, 5, 30),
		CameraZ:          fake.Float(2, -50, 50),
		ThumbnailPath:    fmt.Sprintf("machines/thumbnails/lanjutkan.png"),
		ModelPath:        fmt.Sprintf("machines/models/chiller-1.glb"),
		MachineDocuments: nil,
	}

	if err := gormDatabase.Model(&entity.Machine{}).Create(&machine).Error; err != nil {
		return err
	}

	return nil
}
