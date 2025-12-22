package role

import (
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
	auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		roleRepository:       roleRepository,
		validatorService:     validatorService,
		dbConnection:         dbConnection,
		viperConfig:          viperConfig,
		permissionRepository: permissionRepository,
		auditLogService:      auditLogService,
	}
}

func (roleService *ServiceImpl) FindAll(ginContext *gin.Context) []*model.RoleResponse {
	var roleResponsesRequest []*model.RoleResponse
	backgroundContext := ginContext.Request.Context()
	roleEntities, err := roleService.roleRepository.FindAllCache(backgroundContext)
	helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.StatusInternalError))
	if roleEntities != nil && len(roleEntities) > 0 {
		roleResponsesRequest = helper.MapEntitiesIntoResponsesWithFunc[
			entity.Role,
			*model.RoleResponse,
		](roleEntities, mapper.FuncMapAuditable)
		return roleResponsesRequest
	}
	err = roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {

		roleEntities, err = roleService.roleRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = roleService.roleRepository.SetAllCache(backgroundContext, roleEntities)
		for _, roleEntity := range roleEntities {
			err = roleService.roleRepository.SetByIdCache(ginContext.Request.Context(), roleEntity.Id, &roleEntity)
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.StatusInternalError))
		}
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.StatusInternalError))
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
func (roleService *ServiceImpl) Create(ginContext *gin.Context, createRoleRequest *model.CreateRoleRequest) []*model.RoleResponse {
	var processedId uint64
	valErr := roleService.validatorService.ValidateStruct(createRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *createRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntity := helper.MapCreateRequestIntoEntity[model.CreateRoleRequest, entity.Role](createRoleRequest)
		permissionEntities, err := roleService.permissionRepository.FindBatchById(gormTransaction, createRoleRequest.PermissionIds)
		if len(permissionEntities) != len(createRoleRequest.PermissionIds) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = roleService.roleRepository.Create(gormTransaction, roleEntity)
		processedId = roleEntity.Id
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	backgroundContext := ginContext.Request.Context()

	_ = roleService.roleRepository.DeleteByIdCache(backgroundContext, processedId)
	_ = roleService.roleRepository.DeleteAllCache(backgroundContext)
	return roleService.FindAll(ginContext)
}

func (roleService *ServiceImpl) Update(ginContext *gin.Context, updateRoleRequest *model.UpdateRoleRequest) []*model.RoleResponse {
	var processedId uint64

	valErr := roleService.validatorService.ValidateStruct(updateRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *updateRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntity, err := roleService.roleRepository.FindById(gormTransaction, updateRoleRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		helper.MapUpdateRequestIntoEntity[*model.UpdateRoleRequest, entity.Role](updateRoleRequest, roleEntity)
		err = roleService.roleRepository.Update(gormTransaction, roleEntity)
		permissionEntities, err := roleService.permissionRepository.FindBatchById(gormTransaction, updateRoleRequest.PermissionIds)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		if len(permissionEntities) != len(updateRoleRequest.PermissionIds) {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest))
		}
		err = gormTransaction.Model(roleEntity).Association("Permissions").Replace(permissionEntities)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	backgroundContext := ginContext.Request.Context()
	_ = roleService.roleRepository.DeleteByIdCache(backgroundContext, processedId)
	_ = roleService.roleRepository.DeleteAllCache(backgroundContext)
	return roleService.FindAll(ginContext)
}

func (roleService *ServiceImpl) Delete(ginContext *gin.Context, deleteRoleRequest *model.DeleteResourceGeneralRequest) []*model.RoleResponse {
	var processedId uint64
	userClaim := helper.ExtractJwtClaimFromContext(ginContext)
	valErr := roleService.validatorService.ValidateStruct(deleteRoleRequest)
	roleService.validatorService.ParseValidationError(valErr, *deleteRoleRequest)
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		err := roleService.roleRepository.Delete(gormTransaction, deleteRoleRequest.Id)
		roleService.auditLogService.Create(ginContext, &model.CreateAuditLogRequest{
			UserId:      userClaim.Id,
			Action:      model.AUDIT_LOG_DELETE,
			Feature:     model.AUDIT_LOG_FEATURE_ROLE,
			Description: deleteRoleRequest.Reason,
			Before:      nil,
			After:       nil,
			IpAddress:   "",
		})
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		processedId = deleteRoleRequest.Id
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	backgroundContext := ginContext.Request.Context()
	_ = roleService.roleRepository.DeleteByIdCache(backgroundContext, processedId)
	_ = roleService.roleRepository.DeleteAllCache(backgroundContext)
	return roleService.FindAll(ginContext)
}

func (roleService *ServiceImpl) FindById(ginContext *gin.Context, roleId uint64) *model.RoleResponse {
	var roleResponse *model.RoleResponse
	backgroundContext := ginContext.Request.Context()
	err := roleService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		roleEntity, err := roleService.roleRepository.FindRoleCacheById(backgroundContext, roleId)
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.StatusInternalError))
		if roleEntity != nil {
			roleResponse = helper.MapEntityIntoResponse[
				*entity.Role,
				*model.RoleResponse,
			](roleEntity,
				[]string{},
				mapper.FuncMapAuditable)
			return nil
		}
		roleEntity, err = roleService.roleRepository.FindById(gormTransaction, roleId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		roleResponse = helper.MapEntityIntoResponse[
			*entity.Role,
			*model.RoleResponse,
		](roleEntity,
			[]string{},
			mapper.FuncMapAuditable)
		err = roleService.roleRepository.SetByIdCache(backgroundContext, roleId, roleEntity)
		helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusInternalServerError, exception.StatusInternalError))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return roleResponse
}
