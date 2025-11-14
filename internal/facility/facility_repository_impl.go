package facility

import (
	"go-intconnect-api/internal/entity"

	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (facilityRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Facility, error) {
	var facilityEntities []entity.Facility
	err := gormTransaction.Find(&facilityEntities).Error
	return facilityEntities, err
}

func (facilityRepositoryImpl *RepositoryImpl) FindBatchById(gormTransaction *gorm.DB, facilityIds []uint64) ([]entity.Facility, error) {
	var facilityEntities []entity.Facility
	err := gormTransaction.Where("id IN ?", facilityIds).Find(&facilityEntities).Error
	return facilityEntities, err
}

func (facilityRepositoryImpl *RepositoryImpl) FindAllPagination(gormTransaction *gorm.DB, orderClause string, offsetVal, limitPage int, searchQuery string) ([]entity.Facility, int64, error) {
	var facilityEntities []entity.Facility
	var totalItems int64

	if searchQuery != "" {
		// Add search condition
		searchPattern := "%" + searchQuery + "%"
		gormTransaction = gormTransaction.Where("facilityname LIKE ? OR email LIKE ?  OR password = ?", searchPattern, searchPattern, searchPattern)

	}

	// Count total items
	err := gormTransaction.Model(&entity.Facility{}).
		Preload("FacilityGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Order(orderClause).Offset(offsetVal).Limit(limitPage).Find(&facilityEntities).Error
	gormTransaction.Model(&entity.Facility{}).Count(&totalItems)
	return facilityEntities, totalItems, err
}

func (facilityRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, facilityId uint64) (*entity.Facility, error) {
	var facilityEntity entity.Facility
	err := gormTransaction.Model(&entity.Facility{}).
		Preload("FacilityGroup", func(gormTx *gorm.DB) *gorm.DB {
			return gormTx.Select("id, name")
		}).Where("id = ?", facilityId).Find(&facilityEntity).Error

	return &facilityEntity, err
}

func (facilityRepositoryImpl *RepositoryImpl) FindByName(pipelineName string) *entity.Facility {
	//TODO implement me
	panic("implement me")
}

func (facilityRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Facility) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (facilityRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, facilityEntity *entity.Facility) error {
	return gormTransaction.Model(facilityEntity).Save(facilityEntity).Error
}

func (facilityRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(entity.Facility{}).Where("id = ?", id).Delete(entity.Facility{}).Error
}
