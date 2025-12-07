package entity

type SmtpServer struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Host        string    `gorm:"column:host"`
	Port        string    `gorm:"column:port"`
	Username    string    `gorm:"column:username"`
	Password    string    `gorm:"column:password"`
	MailAddress string    `gorm:"column:mail_address"`
	MailName    string    `gorm:"column:mail_name"`
	IsActive    bool      `gorm:"column:is_active"`
	Auditable   Auditable `gorm:"embedded"`
}

func (smtpServerEntity SmtpServer) GetAuditable() *Auditable {
	return &smtpServerEntity.Auditable
}
