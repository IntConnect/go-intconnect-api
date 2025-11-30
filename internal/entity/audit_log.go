package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type AuditLog struct {
	Id              uint64                 `gorm:"column:id;primaryKey;autoIncrement"`
	UserId          uint64                 `gorm:"column:user_id"`
	Action          string                 `gorm:"column:action"`
	Feature         string                 `gorm:"column:feature"`
	Description     string                 `gorm:"column:description"`
	Before          map[string]interface{} `gorm:"-:all"`
	After           map[string]interface{} `gorm:"-:all"`
	BeforeRaw       []byte                 `gorm:"column:before;type:jsonb"`
	AfterRaw        []byte                 `gorm:"column:after;type:jsonb"`
	IpAddress       string                 `gorm:"column:ip_address"`
	User            User                   `gorm:"foreignKey:UserId;references:Id"`
	SimpleAuditable SimpleAuditable        `gorm:"embedded"`
}

func (auditLogEntity *AuditLog) AfterFind(gormTransaction *gorm.DB) (err error) {
	if len(auditLogEntity.BeforeRaw) > 0 {
		err = json.Unmarshal(auditLogEntity.BeforeRaw, &auditLogEntity.Before)
	}
	if len(auditLogEntity.AfterRaw) > 0 {
		err = json.Unmarshal(auditLogEntity.AfterRaw, &auditLogEntity.After)
	}
	return
}
func (auditLogEntity *AuditLog) BeforeSave(gormTransaction *gorm.DB) (err error) {
	if auditLogEntity.Before != nil && len(auditLogEntity.Before) > 0 {
		auditLogEntity.BeforeRaw, err = json.Marshal(auditLogEntity.Before)
	}
	if auditLogEntity.After != nil && len(auditLogEntity.After) > 0 {
		auditLogEntity.AfterRaw, err = json.Marshal(auditLogEntity.After)
	}
	return nil
}

func (auditLogEntity AuditLog) GetSimpleAuditable() *SimpleAuditable {
	return &auditLogEntity.SimpleAuditable

}
