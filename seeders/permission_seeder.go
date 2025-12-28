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

		// Role & User
		{
			Code:        "ROLE_USER_VIEW",
			Name:        "View",
			Category:    "Role & User",
			Description: "Permission to view roles & users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_USER_CREATE",
			Name:        "Create",
			Category:    "Role & User",
			Description: "Permission to create roles & users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_USER_EDIT",
			Name:        "Update",
			Category:    "Role & User",
			Description: "Permission to update roles & users",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "ROLE_USER_DELETE",
			Name:        "Delete",
			Category:    "Role & User",
			Description: "Permission to delete roles & users",
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

		// BREAKDOWN
		{
			Code:        "BREAKDOWN_VIEW",
			Name:        "View",
			Category:    "Breakdown",
			Description: "Permission to view breakdowns",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "BREAKDOWN_CREATE",
			Name:        "Create",
			Category:    "Breakdown",
			Description: "Permission to create breakdowns",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "BREAKDOWN_EDIT",
			Name:        "Update",
			Category:    "Breakdown",
			Description: "Permission to update breakdowns",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "BREAKDOWN_DELETE",
			Name:        "Delete",
			Category:    "Breakdown",
			Description: "Permission to delete breakdowns",
			Auditable:   entity.NewAuditable("System"),
		},

		// LOG ALARM
		{
			Code:        "LOG_ALARM_VIEW",
			Name:        "View",
			Category:    "Log Alarm",
			Description: "Permission to view log alarm",
			Auditable:   entity.NewAuditable("System"),
		},

		// Report Document Template
		{
			Code:        "REPORT_DOCUMENT_TEMPLATE_VIEW",
			Name:        "View",
			Category:    "Report Document Template",
			Description: "Permission to view report document templates",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "REPORT_DOCUMENT_TEMPLATE_CREATE",
			Name:        "Create",
			Category:    "Report Document Template",
			Description: "Permission to create report document templates",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "REPORT_DOCUMENT_TEMPLATE_EDIT",
			Name:        "Update",
			Category:    "Report Document Template",
			Description: "Permission to update report document templates",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "REPORT_DOCUMENT_TEMPLATE_DELETE",
			Name:        "Delete",
			Category:    "Report Document Template",
			Description: "Permission to delete report document templates",
			Auditable:   entity.NewAuditable("System"),
		},

		// Generate Report
		{
			Code:        "GENERATE_REPORT_VIEW",
			Name:        "View",
			Category:    "Generate Report",
			Description: "Permission to view generate reports",
			Auditable:   entity.NewAuditable("System"),
		},

		// Permission
		{
			Code:        "PERMISSION_VIEW",
			Name:        "View",
			Category:    "Permission",
			Description: "Permission to view permissions",
			Auditable:   entity.NewAuditable("System"),
		},

		// Audit Log
		{
			Code:        "AUDIT_LOG_VIEW",
			Name:        "View",
			Category:    "Audit Log",
			Description: "Permission to view audit logs",
			Auditable:   entity.NewAuditable("System"),
		},

		// Check Sheet Document Template
		{
			Code:        "CHECK_SHEET_DOCUMENT_TEMPLATE_VIEW",
			Name:        "View",
			Category:    "Check Sheet Document Template",
			Description: "Permission to view check sheet document templates",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "CHECK_SHEET_DOCUMENT_TEMPLATE_CREATE",
			Name:        "Create",
			Category:    "Check Sheet Document Template",
			Description: "Permission to create check sheet document templates",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "CHECK_SHEET_DOCUMENT_TEMPLATE_EDIT",
			Name:        "Update",
			Category:    "Check Sheet Document Template",
			Description: "Permission to update check sheet document templates",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "CHECK_SHEET_DOCUMENT_TEMPLATE_DELETE",
			Name:        "Delete",
			Category:    "Check Sheet Document Template",
			Description: "Permission to delete check sheet document templates",
			Auditable:   entity.NewAuditable("System"),
		},

		// System Setting
		{
			Code:        "SYSTEM_SETTING_VIEW",
			Name:        "View",
			Category:    "System Setting",
			Description: "Permission to view system setting",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "SYSTEM_SETTING_MANAGE",
			Name:        "Manage",
			Category:    "System Setting",
			Description: "Permission to manage system setting",
			Auditable:   entity.NewAuditable("System"),
		},

		// Check Sheet
		{
			Code:        "CHECK_SHEET_VIEW",
			Name:        "View",
			Category:    "Check Sheet",
			Description: "Permission to view check sheet",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "CHECK_SHEET_CREATE",
			Name:        "Create",
			Category:    "Check Sheet",
			Description: "Permission to create check sheet",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "CHECK_SHEET_EDIT",
			Name:        "Update",
			Category:    "Check Sheet",
			Description: "Permission to update check sheet",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "CHECK_SHEET_DELETE",
			Name:        "Delete",
			Category:    "Check Sheet",
			Description: "Permission to delete check sheet",
			Auditable:   entity.NewAuditable("System"),
		},

		// SMTP Server
		{
			Code:        "SMTP_SERVER_VIEW",
			Name:        "View",
			Category:    "SMTP Server",
			Description: "Permission to view SMTP servers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "SMTP_SERVER_CREATE",
			Name:        "Create",
			Category:    "SMTP Server",
			Description: "Permission to create SMTP servers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "SMTP_SERVER_EDIT",
			Name:        "Update",
			Category:    "SMTP Server",
			Description: "Permission to update SMTP servers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "SMTP_SERVER_DELETE",
			Name:        "Delete",
			Category:    "SMTP Server",
			Description: "Permission to delete SMTP servers",
			Auditable:   entity.NewAuditable("System"),
		},

		// Modbus Server
		{
			Code:        "MODBUS_SERVER_VIEW",
			Name:        "View",
			Category:    "MODBUS SERVERS",
			Description: "Permission to view Modbus servers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MODBUS_SERVER_CREATE",
			Name:        "Create",
			Category:    "MODBUS SERVERS",
			Description: "Permission to create Modbus servers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MODBUS_SERVER_EDIT",
			Name:        "Update",
			Category:    "MODBUS SERVERS",
			Description: "Permission to update Modbus servers",
			Auditable:   entity.NewAuditable("System"),
		},
		{
			Code:        "MODBUS_SERVER_DELETE",
			Name:        "Delete",
			Category:    "MODBUS SERVERS",
			Description: "Permission to delete Modbus servers",
			Auditable:   entity.NewAuditable("System"),
		},
	})

	return nil
}
