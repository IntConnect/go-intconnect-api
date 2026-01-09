package smtp_server

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]*entity.SmtpServer, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]*entity.SmtpServer, int64, error)
	FindById(gormTransaction *gorm.DB, smtpServerId uint64) (*entity.SmtpServer, error)
	Create(gormTransaction *gorm.DB, smtpServerEntity *entity.SmtpServer) error
	Update(gormTransaction *gorm.DB, smtpServerEntity *entity.SmtpServer) error
	Delete(gormTransaction *gorm.DB, smtpServerId uint64) error
}
