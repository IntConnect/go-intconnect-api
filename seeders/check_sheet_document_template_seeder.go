package seeders

import (
	"go-intconnect-api/internal/entity"
	"time"

	"gorm.io/gorm"
)

type CheckSheetDocumentTemplateSeeder struct{}

func (checkSheetDocumentTemplateSeeder *CheckSheetDocumentTemplateSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Exec("TRUNCATE TABLE check_sheet_document_templates RESTART IDENTITY CASCADE")

	gormDatabase.Model(&entity.CheckSheetDocumentTemplate{}).Create(&entity.CheckSheetDocumentTemplate{
		MachineId:      1,
		Name:           "Test",
		No:             "Test",
		Description:    "Test",
		Category:       "Cleaning",
		Interval:       1,
		IntervalType:   "Hour",
		RotationType:   "Daily",
		RevisionNumber: 0,
		EffectiveDate:  time.Now(),
		StartingHour:   "12:00:00",
	})

	return nil
}
