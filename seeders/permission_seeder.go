package seeders

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type PermissionSeeder struct{}

func (permissionSeeder *PermissionSeeder) Run(gormDatabase *gorm.DB) error {
	gormDatabase.Model(&entity.Permission{}).Create([]entity.Permission{
		// Machine
		{
			Code:        "MACHINE_VIEW",
			Name:        "View",
			Category:    "Machine",
			Description: "Permission to view machine details",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MACHINE_CREATE",
			Name:        "Create",
			Category:    "Machine",
			Description: "Permission to create machine records",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MACHINE_EDIT",
			Name:        "Update",
			Category:    "Machine",
			Description: "Permission to update machine records",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MACHINE_DELETE",
			Name:        "Delete",
			Category:    "Machine",
			Description: "Permission to delete machine records",
			Auditable:   entity.NewAuditable("System"),
		},

		// MQTT Topic
		{
			Code:        "MQTT_TOPIC_VIEW",
			Name:        "View",
			Category:    "MQTT Topic",
			Description: "Permission to view MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_TOPIC_CREATE",
			Name:        "Create",
			Category:    "MQTT Topic",
			Description: "Permission to create MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_TOPIC_EDIT",
			Name:        "Update",
			Category:    "MQTT Topic",
			Description: "Permission to update MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_TOPIC_DELETE",
			Name:        "Delete",
			Category:    "MQTT Topic",
			Description: "Permission to delete MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},

		// MQTT Broker
		{
			Code:        "MQTT_BROKER_VIEW",
			Name:        "View",
			Category:    "MQTT Broker",
			Description: "Permission to view MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_BROKER_CREATE",
			Name:        "Create",
			Category:    "MQTT Broker",
			Description: "Permission to create MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_BROKER_EDIT",
			Name:        "Update",
			Category:    "MQTT Broker",
			Description: "Permission to update MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_BROKER_DELETE",
			Name:        "Delete",
			Category:    "MQTT Broker",
			Description: "Permission to delete MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},

		// Parameter
		{
			Code:        "PARAMETER_VIEW",
			Name:        "View",
			Category:    "Parameter",
			Description: "Permission to view parameters",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "PARAMETER_CREATE",
			Name:        "Create",
			Category:    "Parameter",
			Description: "Permission to create parameters",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "PARAMETER_EDIT",
			Name:        "Update",
			Category:    "Parameter",
			Description: "Permission to update parameters",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "PARAMETER_DELETE",
			Name:        "Delete",
			Category:    "Parameter",
			Description: "Permission to delete parameters",
			Auditable:   entity.NewAuditable("System"),
		},

		// Role
		{
			Code:        "ROLE_VIEW",
			Name:        "View",
			Category:    "Role",
			Description: "Permission to view roles",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_CREATE",
			Name:        "Create",
			Category:    "Role",
			Description: "Permission to create roles",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_EDIT",
			Name:        "Update",
			Category:    "Role",
			Description: "Permission to update roles",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_DELETE",
			Name:        "Delete",
			Category:    "Role",
			Description: "Permission to delete roles",
			Auditable:   entity.NewAuditable("System"),
		},

		// User
		{
			Code:        "USER_VIEW",
			Name:        "View",
			Category:    "User",
			Description: "Permission to view users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "USER_CREATE",
			Name:        "Create",
			Category:    "User",
			Description: "Permission to create users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "USER_EDIT",
			Name:        "Update",
			Category:    "User",
			Description: "Permission to update users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "USER_DELETE",
			Name:        "Delete",
			Category:    "User",
			Description: "Permission to delete users",
			Auditable:   entity.NewAuditable("System"),
		},

		// Facility
		{
			Code:        "FACILITY_VIEW",
			Name:        "View",
			Category:    "Facility",
			Description: "Permission to view facilities",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "FACILITY_CREATE",
			Name:        "Create",
			Category:    "Facility",
			Description: "Permission to create facilities",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "FACILITY_EDIT",
			Name:        "Update",
			Category:    "Facility",
			Description: "Permission to update facilities",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "FACILITY_DELETE",
			Name:        "Delete",
			Category:    "Facility",
			Description: "Permission to delete facilities",
			Auditable:   entity.NewAuditable("System"),
		},
	})

	return nil
}
