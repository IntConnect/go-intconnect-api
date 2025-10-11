package protocol_configuration

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(gormTransaction *gorm.DB) ([]entity.ProtocolConfiguration, error)
	FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.ProtocolConfiguration, int64, error)
	FindById(gormTransaction *gorm.DB, protocolConfigurationId uint64) (*entity.ProtocolConfiguration, error)
	FindByName(protocolConfigurationName string) *entity.ProtocolConfiguration
	Create(gormTransaction *gorm.DB, protocolConfigurationEntity *entity.ProtocolConfiguration) error
	Update(gormTransaction *gorm.DB, protocolConfigurationEntity *entity.ProtocolConfiguration) error
	Delete(gormTransaction *gorm.DB, protocolConfigurationId uint64) error
}
