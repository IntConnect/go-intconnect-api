package role

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	roleRepository   Repository
	validatorService validator.Service
	dbConnection     *gorm.DB
	viperConfig      *viper.Viper
}

func NewService(roleRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		roleRepository:   roleRepository,
		validatorService: validatorService,
		dbConnection:     dbConnection,
		viperConfig:      viperConfig,
	}
}

func (roleService *ServiceImpl) FindAll() []*model.RoleResponse {
	var roleResponsesRequest []*model.RoleResponse
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntities, err := roleService.roleRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		roleResponsesRequest = helper.MapEntitiesIntoResponsesWithFunc[
			entity.Role,
			*model.RoleResponse,
		](
			roleEntities,
			mapper.FuncMapAuditable[entity.Role, *model.RoleResponse],
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return roleResponsesRequest
}

func (roleService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) model.PaginationResponse[*model.RoleResponse] {
	paginationResp := model.PaginationResponse[*model.RoleResponse]{}
	offsetVal := (paginationReq.Page - 1) * paginationReq.Size
	orderClause := paginationReq.Sort
	if paginationReq.Order != "" {
		orderClause += " " + paginationReq.Order
	}
	var allRole []*model.RoleResponse
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntities, totalItems, err := roleService.roleRepository.FindAllPagination(gormTransaction, orderClause, offsetVal, paginationReq.Size, paginationReq.SearchQuery)
		totalPages := int(math.Ceil(float64(totalItems) / float64(paginationReq.Size)))
		allRole = helper.MapEntitiesIntoResponses[entity.Role, model.RoleResponse](roleEntities)
		paginationResp = model.PaginationResponse[*model.RoleResponse]{
			Data:        allRole,
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

// Create - Membuat role baru
func (roleService *ServiceImpl) Create(ginContext *gin.Context, createRoleRequest *model.CreateRoleRequest) {
	valErr := roleService.validatorService.ValidateStruct(createRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *createRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntity := helper.MapCreateRequestIntoEntity[model.CreateRoleRequest, entity.Role](createRoleRequest)
		err := roleService.roleRepository.Create(gormTransaction, roleEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (roleService *ServiceImpl) Update(ginContext *gin.Context, updateRoleRequest *model.UpdateRoleRequest) {
	valErr := roleService.validatorService.ValidateStruct(updateRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *updateRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		role, err := roleService.roleRepository.FindById(gormTransaction, updateRoleRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity[*model.UpdateRoleRequest, entity.Role](updateRoleRequest, role)
		err = roleService.roleRepository.Update(gormTransaction, role)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (roleService *ServiceImpl) Delete(ginContext *gin.Context, deleteRoleRequest *model.DeleteResourceGeneralRequest) {
	valErr := roleService.validatorService.ValidateStruct(deleteRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *deleteRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := roleService.roleRepository.Delete(gormTransaction, deleteRoleRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
