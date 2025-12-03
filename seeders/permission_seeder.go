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
			Name:        "View Machines",
			Category:    "Machine",
			Description: "Permission to view machine details",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MACHINE_CREATE",
			Name:        "Create Machine",
			Category:    "Machine",
			Description: "Permission to add new machine records",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MACHINE_EDIT",
			Name:        "Edit Machine",
			Category:    "Machine",
			Description: "Permission to update existing machine records",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MACHINE_DELETE",
			Name:        "Delete Machine",
			Category:    "Machine",
			Description: "Permission to delete machine records",
			Auditable:   entity.NewAuditable("System"),
		},

		// MQTT Topic
		{
			Code:        "MQTT_TOPIC_VIEW",
			Name:        "View MQTT Topics",
			Category:    "MQTT Topic",
			Description: "Permission to view MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_TOPIC_CREATE",
			Name:        "Create MQTT Topic",
			Category:    "MQTT Topic",
			Description: "Permission to add new MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_TOPIC_EDIT",
			Name:        "Edit MQTT Topic",
			Category:    "MQTT Topic",
			Description: "Permission to update MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_TOPIC_DELETE",
			Name:        "Delete MQTT Topic",
			Category:    "MQTT Topic",
			Description: "Permission to delete MQTT topics",
			Auditable:   entity.NewAuditable("System"),
		},

		// MQTT Broker
		{
			Code:        "MQTT_BROKER_VIEW",
			Name:        "View MQTT Brokers",
			Category:    "MQTT Broker",
			Description: "Permission to view MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_BROKER_CREATE",
			Name:        "Create MQTT Broker",
			Category:    "MQTT Broker",
			Description: "Permission to add new MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_BROKER_EDIT",
			Name:        "Edit MQTT Broker",
			Category:    "MQTT Broker",
			Description: "Permission to update MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MQTT_BROKER_DELETE",
			Name:        "Delete MQTT Broker",
			Category:    "MQTT Broker",
			Description: "Permission to delete MQTT brokers",
			Auditable:   entity.NewAuditable("System"),
		},

		// Parameter
		{
			Code:        "PARAMETER_VIEW",
			Name:        "View Parameters",
			Category:    "Parameter",
			Description: "Permission to view parameters",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "PARAMETER_CREATE",
			Name:        "Create Parameter",
			Category:    "Parameter",
			Description: "Permission to add new parameters",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "PARAMETER_EDIT",
			Name:        "Edit Parameter",
			Category:    "Parameter",
			Description: "Permission to update parameters",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "PARAMETER_DELETE",
			Name:        "Delete Parameter",
			Category:    "Parameter",
			Description: "Permission to delete parameters",
			Auditable:   entity.NewAuditable("System"),
		},

		// Role
		{
			Code:        "ROLE_VIEW",
			Name:        "View Roles",
			Category:    "Role",
			Description: "Permission to view roles",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_CREATE",
			Name:        "Create Role",
			Category:    "Role",
			Description: "Permission to add new roles",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_EDIT",
			Name:        "Edit Role",
			Category:    "Role",
			Description: "Permission to update roles",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_DELETE",
			Name:        "Delete Role",
			Category:    "Role",
			Description: "Permission to delete roles",
			Auditable:   entity.NewAuditable("System"),
		},

		// User
		{
			Code:        "USER_VIEW",
			Name:        "View Users",
			Category:    "User",
			Description: "Permission to view users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "USER_CREATE",
			Name:        "Create User",
			Category:    "User",
			Description: "Permission to add new users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "USER_EDIT",
			Name:        "Edit User",
			Category:    "User",
			Description: "Permission to update users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "USER_DELETE",
			Name:        "Delete User",
			Category:    "User",
			Description: "Permission to delete users",
			Auditable:   entity.NewAuditable("System"),
		},

		// Facility
		{
			Code:        "FACILITY_VIEW",
			Name:        "View Facilities",
			Category:    "Facility",
			Description: "Permission to view facilities",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "FACILITY_CREATE",
			Name:        "Create Facility",
			Category:    "Facility",
			Description: "Permission to add new facilities",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "FACILITY_EDIT",
			Name:        "Edit Facility",
			Category:    "Facility",
			Description: "Permission to update facilities",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "FACILITY_DELETE",
			Name:        "Delete Facility",
			Category:    "Facility",
			Description: "Permission to delete facilities",
			Auditable:   entity.NewAuditable("System"),
		},
	})

	return nil
}
