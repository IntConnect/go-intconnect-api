package seeders

import (
	"go-intconnect-api/internal/entity"
	"time"

	"gorm.io/gorm"
)

type CheckSheetSeeder struct{}

func (checkSheetSeeder *CheckSheetSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Exec("TRUNCATE TABLE check_sheets RESTART IDENTITY CASCADE")

	gormDatabase.Model(&entity.CheckSheet{}).Create(&entity.CheckSheet{
		CheckSheetDocumentTemplateId: 1,
		ReportedBy:                   1,
		VerifiedBy:                   nil,
		Timestamp:                    time.Now(),
		Note:                         "",
		Status:                       "Draft",
		CheckSheetCheckPoint:         nil,
		Auditable:                    entity.NewAuditable("SEEDER"),
	})

	return nil
}
