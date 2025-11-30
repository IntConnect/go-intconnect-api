package role

import (
	"fmt"
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/permission"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	roleRepository       Repository
	permissionRepository permission.Repository
	auditLogService      auditLog.Service
	validatorService     validator.Service
	dbConnection         *gorm.DB
	viperConfig          *viper.Viper
}

func NewService(roleRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, permissionRepository permission.Repository,
	auditLogService auditLog.Service) *ServiceImpl {
	return &ServiceImpl{
		roleRepository:       roleRepository,
		validatorService:     validatorService,
		dbConnection:         dbConnection,
		viperConfig:          viperConfig,
		permissionRepository: permissionRepository,
		auditLogService:      auditLogService,
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

// Create - Membuat role baru
func (roleService *ServiceImpl) Create(ginContext *gin.Context, createRoleRequest *model.CreateRoleRequest) {
	valErr := roleService.validatorService.ValidateStruct(createRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *createRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntity := helper.MapCreateRequestIntoEntity[model.CreateRoleRequest, entity.Role](createRoleRequest)
		permissionIds, err := roleService.permissionRepository.FindBatchById(gormTransaction, createRoleRequest.PermissionIds)
		if len(permissionIds) != len(createRoleRequest.PermissionIds) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = roleService.roleRepository.Create(gormTransaction, roleEntity)
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
	userClaim := helper.ExtractJwtClaimFromContext(ginContext)
	fmt.Println(userClaim)
	valErr := roleService.validatorService.ValidateStruct(deleteRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *deleteRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := roleService.roleRepository.Delete(gormTransaction, deleteRoleRequest.Id)
		roleService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      0,
			Action:      "",
			Feature:     "",
			Description: "",
			Before:      nil,
			After:       nil,
			IpAddress:   "",
		})
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
