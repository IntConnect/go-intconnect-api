package permission

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"math"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	permissionRepository Repository
	validatorService     validator.Service
	dbConnection         *gorm.DB
	viperConfig          *viper.Viper
}

func NewService(permissionRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		permissionRepository: permissionRepository,
		validatorService:     validatorService,
		dbConnection:         dbConnection,
		viperConfig:          viperConfig,
	}
}

func (permissionService *ServiceImpl) FindAll() []*model.PermissionResponse {
	var permissionResponsesRequest []*model.PermissionResponse
	err := permissionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		permissionEntities, err := permissionService.permissionRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		permissionResponsesRequest = helper.MapEntitiesIntoResponses[entity.Permission, model.PermissionResponse](permissionEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return permissionResponsesRequest
}

func (permissionService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.PermissionResponse] {
	paginationResp := model.PaginationResponse[*model.PermissionResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allPermission []*model.PermissionResponse
	err := permissionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		permissionEntities, totalItems, err := permissionService.permissionRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allPermission = helper.MapEntitiesIntoResponses[entity.Permission, model.PermissionResponse](permissionEntities)
		paginationResp = model.PaginationResponse[*model.PermissionResponse]{
			Data:        allPermission,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: paginationReq.Page,
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}
