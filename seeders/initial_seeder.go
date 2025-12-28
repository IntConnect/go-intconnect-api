package seeders

import (
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type InitialSeeder struct{}

func (initialSeeder *InitialSeeder) Run(gormDatabase *gorm.DB) error {
	seederNames := []string{"RoleSeeder", "PermissionSeeder", "RolePermissionSeeder", "UserSeeder", "FacilitySeeder", "MachineSeeder"}
	for _, name := range seederNames {
		s, err := GetSeeder(name)
		if err != nil {
			log.Println(err)
			continue
		}
		logrus.Debug("Running seeder:", name)
		if err := s.Run(gormDatabase); err != nil {
			log.Println("Seeder failed:", err)
		} else {
			logrus.Debug("Seeder completed:", name)
		}
	}
	return nil
}
