package user

import (
	"go-intconnect-api/internal/entity"
	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (userRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.User, error) {
	var userEntities []entity.User
	err := gormTransaction.Find(&userEntities).Error
	return userEntities, err
}

func (userRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.User, int64, error) {
	var userEntities []entity.User
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("username LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.User{}).
		Preload("UserGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&userEntities).Error
	gormTransaction.Model(&entity.User{}).Count(&totalItems)
	return userEntities, totalItems, err
}

func (userRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, userId uint64) (*entity.User, error) {
	var userEntity entity.User
	err := gormTransaction.Model(&entity.User{}).
		Preload("UserGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", userId).Find(&userEntity).Error

	return &userEntity, err
}

func (userRepositoryImpl *RepositoryImpl) FindByName(currencyName string) *entity.User {
	//TODO implement me
	panic("implement me")
}

func (userRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, currencyEntity *entity.User) error {
	return gormTransaction.Model(currencyEntity).Create(currencyEntity).Error

}

func (userRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, userEntity *entity.User) error {
	return gormTransaction.Model(userEntity).Save(userEntity).Error
}

func (userRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.User{}).Where("id = ?", id).Delete(entity.User{}).Error
}
