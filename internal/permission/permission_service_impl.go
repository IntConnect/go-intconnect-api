package permission

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"

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

		permissionResponsesRequest = helper.MapEntitiesIntoResponsesWithFunc[
			entity.Permission,
			*model.PermissionResponse,
		](
			permissionEntities,
			mapper.FuncMapAuditable[entity.Permission, *model.PermissionResponse],
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return permissionResponsesRequest
}
func (permissionService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.PermissionResponse] {
	var permissionResponses []*model.PermissionResponse
	var totalItems int64

	paginationQuery := helper.BuildPaginationQuery(paginationReq)

	err := permissionService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		permissionEntities, total, err := permissionService.permissionRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)

		// Check error dari repository
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		// Map entities ke responses
		permissionResponses = helper.MapEntitiesIntoResponsesWithFunc[entity.Permission, *model.PermissionResponse](
			permissionEntities,
			mapper.FuncMapAuditable,
		)

		totalItems = total
		return nil
	})

	// Check error dari transaction
	helper.CheckErrorOperation(err, exception.ParseGormError(err))

	// Construct paginated response
	return helper.NewPaginatedResponseFromResult(
		"Permissions fetched successfully",
		permissionResponses,
		paginationReq,
		totalItems,
	)
}
