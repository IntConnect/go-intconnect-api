package register

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (registerRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]*entity.Register, error) {
	var registerEntities []*entity.Register
	err := gormTransaction.
		Preload("Machine").
		Preload("ModbusServer").
		Find(&registerEntities).Error
	return registerEntities, err
}

func (registerRepositoryImpl *RepositoryImpl) FindAllPagination(
	gormTransaction *gorm.DB,
	orderClause string,
	offsetVal, limitPage int,
	searchQuery string,
) ([]*entity.Register, int64, error) {

	var registerEntities []*entity.Register
	var totalItems int64

	// Base query
	rawQuery := gormTransaction.Model(&entity.Register{})

	// Search
	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		rawQuery = rawQuery.Where("memory_location ILIKE ? OR name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Count first
	if err := rawQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	if err := rawQuery.
		Preload("Machine").
		Preload("ModbusServer").
		Order(orderClause).
		Offset(offsetVal).
		Limit(limitPage).
		Find(&registerEntities).Error; err != nil {
		return nil, 0, err
	}

	return registerEntities, totalItems, nil
}

func (registerRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, registerId uint64) (*entity.Register, error) {
	var registerEntity entity.Register
	err := gormTransaction.Model(&entity.Register{}).
		Preload("ModbusServer").
		Where("id = ?", registerId).Find(&registerEntity).Error

	return &registerEntity, err
}

func (registerRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Register) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (registerRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, registerEntity *entity.Register) error {
	return gormTransaction.Model(registerEntity).Save(registerEntity).Error
}

func (registerRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Register{}).Where("id = ?", id).Delete(&entity.Register{}).Error
}
