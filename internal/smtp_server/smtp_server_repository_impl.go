package smtp_server

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (smtpServerRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.SmtpServer, error) {
	var smtpServerEntities []*entity.SmtpServer
	err := gormTransaction.Find(&smtpServerEntities).Error
	return smtpServerEntities, err
}

func (smtpServerRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.SmtpServer, int64, error) {

	var smtpServerEntities []*entity.SmtpServer
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.SmtpServer{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("host ILIKE ? OR port ILIKE ? OR username ILIKE ? OR password ILIKE ? OR mail_address ILIKE ? mail_name ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&smtpServerEntities).Error; err != nil {
		return nil, 0, err
	}

	return smtpServerEntities, totalItems, nil
}

func (smtpServerRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, smtpServerId uint64) (*entity.SmtpServer, error) {
	var smtpServerEntity entity.SmtpServer
	err := gormTransaction.Model(&entity.SmtpServer{}).
		Where("id = ?", smtpServerId).Find(&smtpServerEntity).Error

	return &smtpServerEntity, err
}

func (smtpServerRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.SmtpServer) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (smtpServerRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, smtpServerEntity *entity.SmtpServer) error {
	return gormTransaction.Model(smtpServerEntity).Save(smtpServerEntity).Error
}

func (smtpServerRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.SmtpServer{}).Where("id = ?", id).Delete(&entity.SmtpServer{}).Error
}
