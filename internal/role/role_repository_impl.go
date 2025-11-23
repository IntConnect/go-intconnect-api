package role

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (roleRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Role, error) {
	var roleEntities []entity.Role
	err := gormTransaction.Find(&roleEntities).Error
	return roleEntities, err
}

func (roleRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, roleId uint64) (*entity.Role, error) {
	var roleEntity entity.Role
	err := gormTransaction.Model(&entity.Role{}).Where("id = ?", roleId).Find(&roleEntity).Error

	return &roleEntity, err
}

func (roleRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Role) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (roleRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, roleEntity *entity.Role) error {
	return gormTransaction.Model(roleEntity).Save(roleEntity).Error
}

func (roleRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Role{}).Where("id = ?", id).Delete(entity.Role{}).Error
}
