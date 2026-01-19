package register

import (
	auditLog "go-intconnect-api/internal/audit_log"
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/internal/validator"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	registerRepository Repository
	auditLogService    auditLog.Service
	validatorService   validator.Service
	dbConnection       *gorm.DB
	viperConfig        *viper.Viper
}

func NewService(registerRepository Repository, validatorService validator.Service, dbConnection *gorm.DB,
	viperConfig *viper.Viper, auditLogService auditLog.Service,
) *ServiceImpl {
	return &ServiceImpl{
		registerRepository: registerRepository,
		auditLogService:    auditLogService,
		validatorService:   validatorService,
		dbConnection:       dbConnection,
		viperConfig:        viperConfig,
	}
}

func (registerService *ServiceImpl) FindAll() []*model.RegisterResponse {
	var registerResponsesRequest []*model.RegisterResponse
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		registerEntities, err := registerService.registerRepository.FindAll(gormTransaction)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		registerResponsesRequest = helper.MapEntitiesIntoResponses[*entity.Register, *model.RegisterResponse](registerEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return registerResponsesRequest
}

func (registerService *ServiceImpl) FindAllPagination(paginationReq *model.PaginationRequest) *model.PaginatedResponse[*model.RegisterResponse] {
	paginationQuery := helper.BuildPaginationQuery(paginationReq)
	var registerResponses []*model.RegisterResponse
	var totalItems int64
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		registerEntities, total, err := registerService.registerRepository.FindAllPagination(
			gormTransaction,
			paginationQuery.OrderClause,
			paginationQuery.Offset,
			paginationQuery.Limit,
			paginationQuery.SearchQuery,
		)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))

		registerResponses = helper.MapEntitiesIntoResponsesWithFunc[*entity.Register, *model.RegisterResponse](
			registerEntities,
			mapper.FuncMapAuditable,
		)
		totalItems = total

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return helper.NewPaginatedResponseFromResult(
		"Register fetched successfully",
		registerResponses,
		paginationReq,
		totalItems,
	)
}

func (registerService *ServiceImpl) FindById(ginContext *gin.Context, registerId uint64) *model.RegisterResponse {
	var registerResponse *model.RegisterResponse
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		registerEntity, err := registerService.registerRepository.FindById(gormTransaction, registerId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		registerResponse = helper.MapEntityIntoResponse[
			*entity.Register,
			*model.RegisterResponse,
		](
			registerEntity,
			mapper.FuncMapAuditable[*entity.Register, *model.RegisterResponse],
		)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return registerResponse
}

func (registerService *ServiceImpl) Create(ginContext *gin.Context, createRegisterRequest *model.CreateRegisterRequest) *model.PaginatedResponse[*model.RegisterResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginationResp *model.PaginatedResponse[*model.RegisterResponse]
	valErr := registerService.validatorService.ValidateStruct(createRegisterRequest)
	registerService.validatorService.ParseValidationError(valErr, *createRegisterRequest)
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		registerEntity := helper.MapCreateRequestIntoEntity[model.CreateRegisterRequest, entity.Register](createRegisterRequest)
		registerEntity.Auditable = entity.NewAuditable(userJwtClaims.Username)
		err := registerService.registerRepository.Create(gormTransaction, registerEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := registerService.auditLogService.Build(
			registerEntity,
			nil,
			nil,
			"",
		)

		err = registerService.auditLogService.Record(ginContext,
			model.AUDIT_LOG_UPDATE,
			model.AUDIT_LOG_FEATURE_REGISTER,
			auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = registerService.FindAllPagination(&paginationRequest)
	return paginationResp
}
func (registerService *ServiceImpl) FindDependency() *model.RegisterDependency {
	var registerDependency = &model.RegisterDependency{}
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var machineEntities []entity.Machine
		var modbusServerEntities []entity.ModbusServer
		err := gormTransaction.Find(&machineEntities).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Find(&modbusServerEntities).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		registerDependency.MachineResponses = helper.MapEntitiesIntoResponses[entity.Machine, model.MachineResponse](machineEntities)
		registerDependency.ModbusServerResponses = helper.MapEntitiesIntoResponses[entity.ModbusServer, model.ModbusServerResponse](modbusServerEntities)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return registerDependency
}

func (registerService *ServiceImpl) Update(ginContext *gin.Context, updateRegisterRequest *model.UpdateRegisterRequest) *model.PaginatedResponse[*model.RegisterResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginationResp *model.PaginatedResponse[*model.RegisterResponse]
	valErr := registerService.validatorService.ValidateStruct(updateRegisterRequest)
	registerService.validatorService.ParseValidationError(valErr, *updateRegisterRequest)
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		registerEntity, err := registerService.registerRepository.FindById(gormTransaction, updateRegisterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		oldRegisterEntity := *registerEntity
		helper.MapUpdateRequestIntoEntity(updateRegisterRequest, registerEntity)
		registerEntity.Auditable = entity.UpdateAuditable(userJwtClaims.Username)
		err = registerService.registerRepository.Update(gormTransaction, registerEntity)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		auditPayload := registerService.auditLogService.Build(
			&oldRegisterEntity,
			registerEntity,
			nil,
			"",
		)

		err = registerService.auditLogService.Record(ginContext,
			model.AUDIT_LOG_UPDATE,
			model.AUDIT_LOG_FEATURE_REGISTER,
			auditPayload)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	paginationRequest := model.NewPaginationRequest()
	paginationResp = registerService.FindAllPagination(&paginationRequest)
	return paginationResp
}

func (registerService *ServiceImpl) Delete(ginContext *gin.Context, deleteRegisterRequest *model.DeleteResourceGeneralRequest) *model.PaginatedResponse[*model.RegisterResponse] {
	userJwtClaims := helper.ExtractJwtClaimFromContext(ginContext)
	var paginationResp *model.PaginatedResponse[*model.RegisterResponse]
	valErr := registerService.validatorService.ValidateStruct(deleteRegisterRequest)
	registerService.validatorService.ParseValidationError(valErr, *deleteRegisterRequest)
	err := registerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		registerEntity, err := registerService.registerRepository.FindById(gormTransaction, deleteRegisterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		registerEntity.Auditable = entity.DeleteAuditable(userJwtClaims.Username)
		err = registerService.registerRepository.Delete(gormTransaction, deleteRegisterRequest.Id)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		paginationRequest := model.NewPaginationRequest()
		paginationResp = registerService.FindAllPagination(&paginationRequest)
		auditPayload := registerService.auditLogService.Build(
			registerEntity,
			nil,
			nil,
			deleteRegisterRequest.Reason,
		)

		err = registerService.auditLogService.Record(ginContext,
			model.AUDIT_LOG_DELETE,
			model.AUDIT_LOG_FEATURE_REGISTER,
			auditPayload)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	paginationRequest := model.NewPaginationRequest()
	paginationResp = registerService.FindAllPagination(&paginationRequest)

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return paginationResp
}
